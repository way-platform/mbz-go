package main

import (
	"context"
	"image/color"
	"net/http"
	"os"

	"charm.land/fang/v2"
	"charm.land/lipgloss/v2"
	"github.com/adrg/xdg"
	"github.com/way-platform/mbz-go"
	"github.com/way-platform/mbz-go/cli"
)

func main() {
	fleetCredPath, _ := xdg.ConfigFile("mbz-go/fleet-credentials.json")
	vspecCredPath, _ := xdg.ConfigFile("mbz-go/vehicle-spec-credentials.json")
	tokenPath, _ := xdg.ConfigFile("mbz-go/token.json")
	var debug bool
	cmd := cli.NewCommand(
		cli.WithFleetCredentialStore(cli.NewFleetCredentialFileStore(fleetCredPath)),
		cli.WithVehicleSpecCredentialStore(cli.NewVehicleSpecCredentialFileStore(vspecCredPath)),
		cli.WithTokenStore(cli.NewFileStore(tokenPath)),
		cli.WithHTTPClient(&http.Client{
			Transport: &mbz.DebugTransport{Enabled: &debug},
		}),
	)
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
	); err != nil {
		os.Exit(1)
	}
}
