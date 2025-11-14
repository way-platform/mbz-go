package vehiclespecificationfleetv1

//go:generate echo [vehiclespecificationfleetv1] copying original...
//go:generate cp vehicle_specification_fleet.yaml 01-original.yaml

//go:generate echo [vehiclespecificationfleetv1] applying overlay...
//go:generate sh -c "go tool -modfile ../../tools/go.mod openapi-overlay apply overlay.yaml 01-original.yaml > 02-overlayed.yaml"

//go:generate echo [vehiclespecificationfleetv1] generating code...
//go:generate go tool -modfile ../../tools/go.mod oapi-codegen -config config.yaml 02-overlayed.yaml
