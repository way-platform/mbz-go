package mbz_test

import (
	"context"
	"fmt"
	"os"

	"github.com/way-platform/mbz-go"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
)

func ExampleClient() {
	ctx := context.Background()
	// Create a Mercedes-Benz API client.
	client, err := mbz.NewClient(
		ctx,
		mbz.WithRegion(mbz.RegionECE),
		mbz.WithClientID(os.Getenv("MBZ_CLIENT_ID")),
		mbz.WithClientSecret(os.Getenv("MBZ_CLIENT_SECRET")),
	)
	if err != nil {
		panic(err)
	}
	// List vehicles in the account.
	response, err := client.ListVehicles(ctx, &fleetv1.ListVehiclesRequest{})
	if err != nil {
		panic(err)
	}
	for _, vehicle := range response.GetVehicles() {
		fmt.Println(vehicle.GetVin())
	}
	// For all available methods, see the API documentation.
}
