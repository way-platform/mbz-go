package mbz

import (
	"github.com/way-platform/mbz-go/api/vehiclesv1"
	mbzv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mbz/v1"
)

func compatibilityResponseToProto(
	vin string,
	resp *vehiclesv1.CompatibilityResponse,
) *mbzv1.VehicleCompatibility {
	msg := &mbzv1.VehicleCompatibility{}
	msg.SetVin(vin)
	msg.SetVehicleType(resp.VehicleType)
	msg.SetVehicleProvidesConnectivity(resp.VehicleProvidesConnectivity)
	services := make([]*mbzv1.VehicleCompatibility_Service, len(resp.Services))
	for i, s := range resp.Services {
		svc := &mbzv1.VehicleCompatibility_Service{}
		svc.SetAvailable(s.Available)
		svc.SetServiceId(s.ServiceID)
		svc.SetServiceName(s.ServiceName)
		services[i] = svc
	}
	msg.SetServices(services)
	return msg
}
