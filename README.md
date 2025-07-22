# Mercedes-Benz Go

A Go SDK and CLI tool for the Mercedes-Benz [Connect Your Fleet APIs](https://developer.mercedes-benz.com/products/connect_your_fleet).

## SDK

### Features

* Support for the [Vehicle Management API](https://developer.mercedes-benz.com/products/connect_your_fleet/specifications/vehicle_management_api)
* Support for the [Services Catalog API](https://developer.mercedes-benz.com/products/connect_your_fleet/specifications/services_catalog_api)

### Installing

```bash
$ go get github.com/way-platform/mbz-go
```

### Using

```go
ctx := context.Background()
// Create a Mercedes-Benz API client.
client, err := mbz.NewClient(
    mbz.WithRegion(mbz.RegionECE),
    mbz.WithOAuth2(
        os.Getenv("MBZ_CLIENT_ID"),
        os.Getenv("MBZ_CLIENT_SECRET"),
    ),
)
if err != nil {
    panic(err)
}
// List vehicles in the account.
response, err := client.ListVehicles(ctx, &mbz.ListVehiclesRequest{})
if err != nil {
    panic(err)
}
for _, vehicle := range response.Vehicles {
    fmt.Println(vehicle.VIN)
}
// For all available methods, see the API documentation.
```

### Developing

#### Building

The project is built using [Mage](https://magefile.org), see
[magefile.go](./magefile.go).

```bash
$ go tool mage build
```

For all available build tasks, see:

```bash
$ go tool mage
```

## CLI tool

The `mbz` CLI tool enables interaction with the APIs from the command line.

<img src="docs/cli.gif" />

### Installing

```bash
go install github.com/way-platform/mbz-go/cmd/mbz@latest
```

Prebuilt binaries for Linux, Windows, and Mac are available from the [Releases](https://github.com/way-platform/mbz-go/releases).
