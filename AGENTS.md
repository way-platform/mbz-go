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
