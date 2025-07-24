package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoparse" //nolint:staticcheck
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type asyncAPISchema struct {
	Components struct {
		Schemas map[string]struct {
			Description string `yaml:"description"`
			Properties  struct {
				Type struct {
					Const string `yaml:"const"`
				} `yaml:"type"`
				Data struct {
					Properties struct {
						Signals struct {
							Items []struct {
								Ref string `yaml:"$ref"`
							} `yaml:"items"`
						} `yaml:"signals"`
					} `yaml:"properties"`
				} `yaml:"data"`
			} `yaml:"properties"`
		} `yaml:"schemas"`
	} `yaml:"components"`
}

func main() {
	const (
		signalProtoPath      = "wayplatform/mbz/v1/signal.proto"
		annotationsProtoPath = "wayplatform/mbz/v1/annotations.proto"
		signalTypeProtoPath  = "wayplatform/mbz/v1/signal_type.proto"
	)
	cmd := &cobra.Command{
		Use:   "update-signal-proto",
		Short: "Update the signal.proto file with the latest signal names from the AsyncAPI file",
	}
	asyncAPIFile := cmd.Flags().String("async-api-file", "", "The path to the input AsyncAPI file")
	_ = cmd.MarkFlagFilename("async-api-file", "yaml")
	_ = cmd.MarkFlagRequired("async-api-file")
	protoDir := cmd.Flags().String("proto-dir", "", "The path to the proto directory")
	_ = cmd.MarkFlagDirname("proto-dir")
	_ = cmd.MarkFlagRequired("proto-dir")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		asyncAPIFileData, err := os.ReadFile(*asyncAPIFile)
		if err != nil {
			return fmt.Errorf("failed to read asyncAPI file: %w", err)
		}
		var schema asyncAPISchema
		if err := yaml.Unmarshal(asyncAPIFileData, &schema); err != nil {
			return fmt.Errorf("failed to unmarshal asyncAPI file: %w", err)
		}
		signalProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal.proto file: %w", err)
		}
		annotationsProtoData, err := os.ReadFile(filepath.Join(*protoDir, annotationsProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read annotations.proto file: %w", err)
		}
		signalTypeProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalTypeProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal.proto file: %w", err)
		}
		fileDescriptors, err := protoparse.Parser{
			Accessor: protoparse.FileContentsFromMap(map[string]string{
				signalProtoPath:      string(signalProtoData),
				annotationsProtoPath: string(annotationsProtoData),
				signalTypeProtoPath:  string(signalTypeProtoData),
			}),
		}.ParseFiles(signalProtoPath, annotationsProtoPath, signalTypeProtoPath)
		if err != nil {
			return fmt.Errorf("failed to parse existing proto files: %w", err)
		}
		fb := builder.NewFile(signalProtoPath)
		fb.SetPackageName("wayplatform.mbz.v1")
		fb.IsProto3 = true
		fb.AddImportedDependency(fileDescriptors[1])
		enumBuilder := builder.NewEnum("SignalName")
		enumBuilder.SetComments(builder.Comments{
			LeadingComment: "SignalName is an enum of all known Mercedes-Benz Kafka message signal names.",
		})
		existingEnum := fileDescriptors[0].FindEnum("wayplatform.mbz.v1.SignalName")
		if existingEnum == nil {
			return fmt.Errorf("SignalName enum not found in signal.proto")
		}
		addedSignalNames := make(map[string]bool)
		enumBuilder.AddValue(
			builder.NewEnumValue("SIGNAL_NAME_UNSPECIFIED").
				SetNumber(0).
				SetComments(builder.Comments{
					LeadingComment: "Default value. This value is unused.",
				}),
		)
		addedSignalNames[""] = true
		maxNumber := int32(0)
		for _, enumValue := range existingEnum.GetValues() {
			if enumValue.GetName() == "SIGNAL_NAME_UNSPECIFIED" {
				continue
			}
			var signalName, signalType string
			if options := enumValue.GetEnumValueOptions(); options != nil {
				for _, uninterpreted := range options.GetUninterpretedOption() {
					if len(uninterpreted.GetName()) == 0 {
						continue
					}
					if !uninterpreted.GetName()[0].GetIsExtension() {
						continue
					}
					switch uninterpreted.GetName()[0].GetNamePart() {
					case "signal_name":
						signalName = string(uninterpreted.GetStringValue())
					case "signal_type":
						signalType = string(uninterpreted.GetIdentifierValue())
					}
				}
			}
			if addedSignalNames[signalName] || signalType == "" {
				continue
			}
			enumBuilder.AddValue(
				newEnumValue(
					enumValue.GetName(),
					enumValue.GetNumber(),
					signalName,
					signalType,
					enumValue.GetSourceInfo().GetLeadingComments(),
				),
			)
			addedSignalNames[signalName] = true
			if enumValue.GetNumber() > maxNumber {
				maxNumber = enumValue.GetNumber()
			}
		}
		for _, item := range schema.Components.Schemas["vehiclesignalsevent"].Properties.Data.Properties.Signals.Items {
			signalName := strings.TrimPrefix(item.Ref, "#/components/schemas/")
			enumName := strings.ReplaceAll(strings.ToUpper(signalName), ".", "_")
			if addedSignalNames[signalName] {
				continue
			}
			maxNumber++
			signalSchema, ok := schema.Components.Schemas[signalName]
			if !ok {
				return fmt.Errorf("signal schema not found: %s", signalName)
			}
			signalType := signalSchema.Properties.Type.Const
			if signalType == "" {
				return fmt.Errorf("signal type not found: %s", signalName)
			}
			comment := strings.TrimSpace(signalSchema.Description)
			comment = strings.TrimPrefix(comment, signalName)
			comment = strings.TrimSpace(comment)
			enumBuilder.AddValue(newEnumValue(enumName, maxNumber, signalName, signalType, comment))
			addedSignalNames[signalName] = true
		}
		fb.AddEnum(enumBuilder)
		fileDesc, err := fb.Build()
		if err != nil {
			return fmt.Errorf("failed to build file descriptor: %w", err)
		}
		var pr protoprint.Printer
		str, err := pr.PrintProtoToString(fileDesc)
		if err != nil {
			return fmt.Errorf("failed to print proto: %w", err)
		}
		return os.WriteFile(filepath.Join(*protoDir, signalProtoPath), []byte(str), 0o644)
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newEnumValue(name string, number int32, stringName, signalType, comment string) *builder.EnumValueBuilder {
	result := builder.NewEnumValue(name).SetNumber(number)
	result.SetComments(builder.Comments{
		LeadingComment: comment,
	})
	result.SetOptions(&descriptorpb.EnumValueOptions{
		UninterpretedOption: []*descriptorpb.UninterpretedOption{
			{
				Name: []*descriptorpb.UninterpretedOption_NamePart{
					{
						NamePart:    proto.String("signal_name"),
						IsExtension: proto.Bool(true),
					},
				},
				StringValue: []byte(stringName),
			},
			{
				Name: []*descriptorpb.UninterpretedOption_NamePart{
					{
						NamePart:    proto.String("signal_type"),
						IsExtension: proto.Bool(true),
					},
				},
				IdentifierValue: proto.String(strings.ToUpper(signalType)),
			},
		},
	})
	return result
}
