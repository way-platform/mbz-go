package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclespecificationv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

// GetVehicleSpecificationRequest is the request for [Client.GetVehicleSpecification].
type GetVehicleSpecificationRequest struct {
	// VehicleID is the FIN or VIN of the vehicle (17 characters).
	VehicleID string

	// Locale is the market locale for country and language (e.g., "en_US", "de_DE").
	// Required parameter.
	Locale string
}

// GetVehicleSpecificationResponse is the response for [Client.GetVehicleSpecification].
type GetVehicleSpecificationResponse struct {
	// VehicleSpecification contains the vehicle specification data.
	VehicleSpecification *mbzv1.VehicleSpecification
}

// GetVehicleSpecification gets the vehicle marketing information for a given vehicle ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleSpecification(
	ctx context.Context,
	request *GetVehicleSpecificationRequest,
	opts ...ClientOption,
) (_ *GetVehicleSpecificationResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle specification: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	if request.VehicleID == "" {
		return nil, fmt.Errorf("vehicle ID is required")
	}
	if request.Locale == "" {
		return nil, fmt.Errorf("locale is required")
	}
	values := url.Values{}
	values.Set("locale", request.Locale)
	requestURL := fmt.Sprintf("https://api.mercedes-benz.com/vehicle_specifications/v1/vehicles/%s?%s",
		request.VehicleID, values.Encode())
	httpRequest, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := c.httpClient(cfg).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	responseData, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var openAPIResp vehiclespecificationv1.VehicleSpecificationResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	proto := vehicleSpecificationToProto(&openAPIResp)
	return &GetVehicleSpecificationResponse{
		VehicleSpecification: proto,
	}, nil
}

func vehicleSpecificationToProto(
	openAPIResp *vehiclespecificationv1.VehicleSpecificationResponse,
) *mbzv1.VehicleSpecification {
	if openAPIResp == nil || openAPIResp.VehicleData == nil {
		return &mbzv1.VehicleSpecification{}
	}
	vehicleData := openAPIResp.VehicleData
	protoSpec := &mbzv1.VehicleSpecification{}
	if vehicleData.ModelName != "" {
		protoSpec.SetModelName(vehicleData.ModelName)
	}
	if vehicleData.ModelYear != "" {
		protoSpec.SetModelYear(vehicleData.ModelYear)
	}
	if vehicleData.Brand != nil && vehicleData.Brand.Text != "" {
		protoSpec.SetBrand(vehicleData.Brand.Text)
	}
	if vehicleData.Fuel != nil && vehicleData.Fuel.Text != "" {
		protoSpec.SetFuelType(vehicleData.Fuel.Text)
	}
	if vehicleData.Emissionstandard != nil && vehicleData.Emissionstandard.Text != "" {
		protoSpec.SetEmissionStandard(vehicleData.Emissionstandard.Text)
	}
	if vehicleData.Weight != nil && vehicleData.Weight.Total != nil {
		protoSpec.SetTotalWeight(*vehicleData.Weight.Total)
	}
	if vehicleData.PrimaryEngine != nil {
		if engine := parseEngine(vehicleData.PrimaryEngine); engine != nil {
			protoSpec.SetPrimaryEngine(engine)
		}
	}
	if vehicleData.SecondaryEngine != nil {
		if engine := parseEngine(vehicleData.SecondaryEngine); engine != nil {
			protoSpec.SetSecondaryEngine(engine)
		}
	}
	return protoSpec
}

func parseEngine(openAPIEngine *vehiclespecificationv1.Engine) *mbzv1.VehicleSpecification_Engine {
	if openAPIEngine == nil {
		return nil
	}
	protoEngine := &mbzv1.VehicleSpecification_Engine{}
	if openAPIEngine.Battery != nil && openAPIEngine.Battery.Capacity != nil {
		protoEngine.SetBatteryCapacityKwh(*openAPIEngine.Battery.Capacity)
	}
	if openAPIEngine.FuelTankCapacity != nil {
		protoEngine.SetFuelTankCapacityL(*openAPIEngine.FuelTankCapacity)
	}
	if openAPIEngine.FuelType != nil && openAPIEngine.FuelType.Text != "" {
		protoEngine.SetFuelType(openAPIEngine.FuelType.Text)
	}
	return protoEngine
}
