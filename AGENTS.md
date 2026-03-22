# Agent Instructions

## Package Manager

Use **Go Modules**: `go mod tidy`, `go test ./...`
Use **Mage**: `go tool mage [target]` (e.g. `go tool mage Build`)

## Key Conventions

- **Testing**: Use standard `testing` and `github.com/google/go-cmp/cmp` **only**. No frameworks (Testify, Ginkgo, etc.).
- **Linting**: Run `GolangCI-Lint` v2. Configure via project-specific `.golangci.yml`.
- **Build**: Use `way-magefile` skill.
- **Encore**: Use `encore-go-*` skills. Encore conventions (e.g., globals) take precedence.

## Local Skills

- **Way Go Style**: Use `.agents/skills/way-go-style/SKILL.md`
- **Way Magefile**: Use `.agents/skills/way-magefile/SKILL.md`
- **Agents.md**: Use `.agents/skills/agents-md/SKILL.md`

## CLI Architecture

The CLI is split into two layers to keep credential storage pluggable:

```
cli/
├── cli.go       # Store interface, Credentials, FileStore, Options
└── command.go   # NewCommand() — full command tree
cmd/mbz/
└── main.go      # Thin wrapper: wires FileStore to XDG paths
```

- `cli.Store` — interface with `Read(any)`, `Write(any)`, `Clear()` methods
- `cli.NewCommand(...Option)` — builds the Cobra command tree; receives stores via functional options (`WithCredentialStore`, `WithTokenStore`)
- `cmd/mbz/main.go` — only wires `FileStore` instances and calls `cli.NewCommand()`

This separation lets consumers embed the CLI in a larger tool or swap the storage backend (e.g. use an in-memory store in tests, or a keychain-backed store) without forking.

### Dual Authentication

The Mercedes-Benz API uses two auth methods depending on the endpoint:

- **OAuth2 client credentials** — for vehicle management, data services, delta push
- **API key** — for vehicle specification and images

Both are stored in the credential store. The OAuth2 token is cached separately in the token store.

### Embedding in a Parent CLI

The CLI can be embedded as a subcommand in a larger tool (e.g. a unified `way` CLI). Key design rules:

- **Never use `cmd.Root()`** — resolves to the parent CLI's root when embedded, breaking flag lookups. Use `cmd.Flags()` instead (works for both persistent and local flags).
- **`WithHTTPClient`** — the parent injects an `*http.Client` via `cli.WithHTTPClient()`. The SDK layers (auth, retry) stack on top of the injected client's transport.
- **`DebugTransport`** — exported in `debug.go` with a lazy `Enabled *bool` field. The parent owns the `--debug` flag and points `Enabled` at the flag variable. The transport checks the pointer at request time, solving the chicken-and-egg problem (transport constructed before flag parsing).

```go
var debug bool
cmd := cli.NewCommand(
    cli.WithCredentialStore(store),
    cli.WithTokenStore(tokenStore),
    cli.WithHTTPClient(&http.Client{
        Transport: &mbz.DebugTransport{
            Enabled: &debug,
            Next:    http.DefaultTransport,
        },
    }),
)
cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
```

### Module Structure

Three separate Go modules prevent Cobra/CLI dependencies from leaking into the SDK library:

```
go.mod          # SDK client library (no cobra, no CLI deps)
cli/go.mod      # CLI commands (depends on root SDK + cobra)
cmd/mbz/go.mod  # Standalone binary (depends on cli module)
```

Consumers who only need the Go client import the root module without pulling in CLI dependencies.

### Conventions

- Subcommands are organized by entity using `cobra.Group`
- Flat command structure: `vehicles`, `services`, not `vehicles list`
