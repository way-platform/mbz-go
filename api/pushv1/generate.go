package pushv1

//go:generate echo [pushv1] copying original...
//go:generate cp asyncapi.yaml 01-original.yaml

//go:generate echo [pushv1] applying overlay...
//go:generate sh -c "go tool -modfile ../../tools/go.mod openapi-overlay apply overlay.yaml 01-original.yaml > 02-overlayed.yaml"
