package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/vehiclespecificationfleetv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

// GetVehicleSpecificationRequest is the request for [Client.GetVehicleSpecification].
type GetVehicleSpecificationRequest struct {
	// VIN is the VIN (or FIN) of the vehicle (17 characters).
	VIN string `json:"vin"`
	// Locale is the market locale.
	Locale string `json:"locale"`
}

// GetVehicleSpecification gets the vehicle marketing information for a given vehicle ID.
// This method requires API key authentication via [WithAPIKey].
func (c *Client) GetVehicleSpecification(
	ctx context.Context,
	request *GetVehicleSpecificationRequest,
	opts ...ClientOption,
) (_ *mbzv1.VehicleSpecification, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: get vehicle specification: %w", err)
		}
	}()
	cfg := c.config.with(opts...)
	if request.VIN == "" {
		return nil, fmt.Errorf("VIN is required")
	}
	values := url.Values{}
	if request.Locale != "" {
		values.Set("locale", request.Locale)
	} else {
		values.Set("locale", string(vehiclespecificationfleetv1.LocalesEnUS))
	}
	// Set all optional parameters to true to maximize data retrieval
	values.Set("additionalSpecs", "true")
	values.Set("optionsNullDescription", "true")
	values.Set("options", "true")
	values.Set("technicalData", "true")
	values.Set("payloadNullValues", "true")
	requestURL := fmt.Sprintf("%s/vehicles/%s", vehiclespecificationfleetv1.BaseURL, request.VIN)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpRequest.URL.RawQuery = values.Encode()
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
	var openAPIResp vehiclespecificationfleetv1.VehicleSpecificationResponse
	if err := json.Unmarshal(responseData, &openAPIResp); err != nil {
		return nil, err
	}
	return vehicleSpecificationToProto(&openAPIResp), nil
}

func vehicleSpecificationToProto(
	openAPIResp *vehiclespecificationfleetv1.VehicleSpecificationResponse,
) *mbzv1.VehicleSpecification {
	if openAPIResp == nil || openAPIResp.VehicleData == nil {
		return &mbzv1.VehicleSpecification{}
	}
	vehicleData := openAPIResp.VehicleData
	protoSpec := &mbzv1.VehicleSpecification{}
	// Use Model field first (contains "Sprinter" in example), fall back to ModelName if Model is empty
	if vehicleData.Model != "" {
		protoSpec.SetModelName(vehicleData.Model)
	} else if vehicleData.ModelName != "" {
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
	if vehicleData.Weight != nil && vehicleData.Weight.Total != nil && *vehicleData.Weight.Total > 0 {
		protoSpec.SetTotalWeightKg(*vehicleData.Weight.Total)
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

func parseEngine(openAPIEngine *vehiclespecificationfleetv1.Engine) *mbzv1.VehicleSpecification_Engine {
	if openAPIEngine == nil {
		return nil
	}
	protoEngine := &mbzv1.VehicleSpecification_Engine{}
	if openAPIEngine.Battery != nil && openAPIEngine.Battery.Capacity != nil && *openAPIEngine.Battery.Capacity > 0 {
		protoEngine.SetBatteryCapacityKwh(*openAPIEngine.Battery.Capacity)
	}
	if openAPIEngine.FuelTankCapacity != nil && *openAPIEngine.FuelTankCapacity > 0 {
		protoEngine.SetFuelTankCapacityL(*openAPIEngine.FuelTankCapacity)
	}
	if openAPIEngine.FuelType != nil && openAPIEngine.FuelType.Text != "" {
		protoEngine.SetFuelType(openAPIEngine.FuelType.Text)
	}
	return protoEngine
}
