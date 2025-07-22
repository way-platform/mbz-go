package mbz_test

import (
	"context"
	"fmt"
	"os"

	"github.com/way-platform/mbz-go"
)

func ExampleClient() {
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
}
