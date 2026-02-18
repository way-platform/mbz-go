package remotemaintenancev3

//go:generate echo [remotemaintenancev3] copying original...
//go:generate cp remote_maintenance_support_api.yaml 01-original.yaml

//go:generate echo [remotemaintenancev3] applying overlay...
//go:generate sh -c "go tool -modfile ../../tools/go.mod openapi-overlay apply overlay.yaml 01-original.yaml > 02-overlayed.yaml"

//go:generate echo [remotemaintenancev3] generating code...
//go:generate go tool -modfile ../../tools/go.mod oapi-codegen -config config.yaml 02-overlayed.yaml
