module github.com/way-platform/mbz-go/cmd/mbz

go 1.24.4

replace github.com/way-platform/mbz-go => ../..

require (
	github.com/adrg/xdg v0.5.3
	github.com/spf13/cobra v1.9.1
	github.com/way-platform/mbz-go v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.30.0
)

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
