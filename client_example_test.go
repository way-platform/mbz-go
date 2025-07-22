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
	fmt.Println(response.Vehicles)

	// List available services with detailed info.
	response2, err := client.ListServices(ctx, &mbz.ListServicesRequest{
		Details: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response2.Services)

	// For all available methods, see the API documentation.
}
