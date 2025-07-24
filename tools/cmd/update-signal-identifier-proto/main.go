package main

import (
	"fmt"
	"iter"
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
		signalIdentifierProtoPath = "wayplatform/mbz/v1/signal_identifier.proto"
		annotationsProtoPath      = "wayplatform/mbz/v1/annotations.proto"
		signalTypeProtoPath       = "wayplatform/mbz/v1/signal_type.proto"
		signalUnitProtoPath       = "wayplatform/mbz/v1/signal_unit.proto"
		signalEnumValueProtoPath  = "wayplatform/mbz/v1/signal_enum_value.proto"
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
		signalIdentifierProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalIdentifierProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal_identifier.proto file: %w", err)
		}
		annotationsProtoData, err := os.ReadFile(filepath.Join(*protoDir, annotationsProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read annotations.proto file: %w", err)
		}
		signalTypeProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalTypeProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal.proto file: %w", err)
		}
		signalUnitProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalUnitProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal_unit.proto file: %w", err)
		}
		signalEnumValueProtoData, err := os.ReadFile(filepath.Join(*protoDir, signalEnumValueProtoPath))
		if err != nil {
			return fmt.Errorf("failed to read signal_enum_value.proto file: %w", err)
		}
		fileDescriptors, err := protoparse.Parser{
			Accessor: protoparse.FileContentsFromMap(map[string]string{
				signalIdentifierProtoPath: string(signalIdentifierProtoData),
				annotationsProtoPath:      string(annotationsProtoData),
				signalTypeProtoPath:       string(signalTypeProtoData),
				signalUnitProtoPath:       string(signalUnitProtoData),
				signalEnumValueProtoPath:  string(signalEnumValueProtoData),
			}),
		}.ParseFiles(signalIdentifierProtoPath, annotationsProtoPath)
		if err != nil {
			return fmt.Errorf("failed to parse existing proto files: %w", err)
		}
		fb := builder.NewFile(signalIdentifierProtoPath)
		fb.SetPackageName("wayplatform.mbz.v1")
		fb.IsProto3 = true
		fb.AddImportedDependency(fileDescriptors[1])
		enumBuilder := builder.NewEnum("SignalIdentifier")
		enumBuilder.SetComments(builder.Comments{
			LeadingComment: "SignalIdentifier is an enum of all known Mercedes-Benz Kafka message signal names.",
		})
		existingEnum := fileDescriptors[0].FindEnum("wayplatform.mbz.v1.SignalIdentifier")
		if existingEnum == nil {
			return fmt.Errorf("SignalIdentifier enum not found in signal.proto")
		}
		addedSignalNames := make(map[string]bool)
		enumBuilder.AddValue(
			builder.NewEnumValue("SIGNAL_IDENTIFIER_UNSPECIFIED").
				SetNumber(0).
				SetComments(builder.Comments{
					LeadingComment: "Default value. This value is unused.",
				}),
		)
		addedSignalNames[""] = true
		maxNumber := int32(0)
		for _, enumValue := range existingEnum.GetValues() {
			if enumValue.GetName() == "SIGNAL_IDENTIFIER_UNSPECIFIED" {
				continue
			}
			var signalName, signalType, signalUnit string
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
					case "signal_unit":
						signalUnit = string(uninterpreted.GetIdentifierValue())
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
					signalUnit,
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
			signalType := strings.ToUpper(signalSchema.Properties.Type.Const)
			if signalType == "" {
				return fmt.Errorf("signal type not found: %s", signalName)
			}
			comment := strings.TrimSpace(signalSchema.Description)
			comment = strings.TrimPrefix(comment, signalName)
			comment = strings.TrimSpace(comment)
			enumBuilder.AddValue(newEnumValue(enumName, maxNumber, signalName, signalType, "", comment))
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
		return os.WriteFile(filepath.Join(*protoDir, signalIdentifierProtoPath), []byte(str), 0o644)
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func tryInferSignalUnit(comment string) string {
	switch {
	case strings.Contains(comment, "(km/h)"):
		return "KILOMETER_PER_HOUR"
	case strings.Contains(comment, "(km)"):
		return "KILOMETER"
	case strings.Contains(comment, "(%)"):
		return "PERCENT"
	case strings.Contains(comment, "(kWh/100km)"):
		return "KILOWATT_HOURS_PER_100KM"
	case strings.Contains(comment, "(kWh)"):
		return "KILOWATT_HOURS"
	case strings.Contains(comment, "(kW)"):
		return "KILOWATT"
	case strings.Contains(comment, "(l/100km)"):
		return "LITER_PER_100KM"
	case strings.Contains(comment, "(V)"):
		return "VOLT"
	case strings.Contains(comment, "(°C)"):
		return "CELSIUS"
	case strings.Contains(comment, "(kPa)"):
		return "KILOPASCAL"
	case strings.Contains(comment, "(min)") || strings.Contains(comment, "in minutes"):
		return "MINUTE"
	case strings.Contains(comment, "in days"):
		return "DAY"
	case strings.Contains(comment, "(ms)") || strings.Contains(comment, "in milliseconds"):
		return "MILLISECOND"
	default:
		return ""
	}
}

func isEnumSignal(signalType string, comment string) bool {
	if signalType != "STRING" {
		return false
	}
	for line := range strings.Lines(comment) {
		if strings.HasPrefix(strings.TrimSpace(line), "* ") {
			return true
		}
	}
	return false
}

func inferEnumSignalValues(comment string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for line := range strings.Lines(comment) {
			if !strings.HasPrefix(strings.TrimSpace(line), "* ") {
				continue
			}
			value, description, _ := strings.Cut(strings.TrimSpace(line[2:]), " = ")
			if !yield(value, description) {
				return
			}
		}
	}
}

func trimEnumSignalValues(comment string) string {
	var result strings.Builder
	result.Grow(len(comment))
	lines := strings.Split(comment, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "* ") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		result.WriteString(line)
		if i != len(lines)-1 {
			result.WriteByte('\n')
		}
	}
	return result.String()
}

func newEnumValue(name string, number int32, stringName, signalType, signalUnit, comment string) *builder.EnumValueBuilder {
	result := builder.NewEnumValue(name).SetNumber(number)
	if isEnumSignal(signalType, comment) {
		signalType = "ENUM"
	}
	options := &descriptorpb.EnumValueOptions{
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
	}
	if signalType == "ENUM" {
		for value, description := range inferEnumSignalValues(comment) {
			option := &descriptorpb.UninterpretedOption{
				Name: []*descriptorpb.UninterpretedOption_NamePart{
					{
						NamePart:    proto.String("signal_values"),
						IsExtension: proto.Bool(true),
					},
				},
			}
			if description != "" {
				option.AggregateValue = proto.String(fmt.Sprintf(`value: "%s", description: "%s"`, value, description))
			} else {
				option.AggregateValue = proto.String(fmt.Sprintf(`value: "%s"`, value))
			}
			options.UninterpretedOption = append(options.UninterpretedOption, option)
		}
		comment = trimEnumSignalValues(comment)
	}
	result.SetComments(builder.Comments{
		LeadingComment: comment,
	})
	if signalUnit == "" {
		signalUnit = tryInferSignalUnit(comment)
	}
	if signalUnit != "" {
		options.UninterpretedOption = append(options.UninterpretedOption, &descriptorpb.UninterpretedOption{
			Name: []*descriptorpb.UninterpretedOption_NamePart{
				{
					NamePart:    proto.String("signal_unit"),
					IsExtension: proto.Bool(true),
				},
			},
			IdentifierValue: proto.String(signalUnit),
		})
	}
	result.SetOptions(options)
	return result
}
