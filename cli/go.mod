module github.com/way-platform/mbz-go/cli

go 1.25.0

toolchain go1.26.0

require (
	github.com/spf13/cobra v1.10.2
	github.com/twmb/franz-go v1.20.7
	github.com/way-platform/mbz-go v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.30.0
	golang.org/x/term v0.41.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

replace github.com/way-platform/mbz-go => ../
