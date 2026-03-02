# AGENTS.md

## Cursor Cloud specific instructions

### Overview
WoWSims TBC is a WoW: The Burning Crusade Classic combat simulator with a Go backend (simulation engine) and a TypeScript/Vite frontend. No external services (databases, Redis, etc.) are required.

### Prerequisites (system-level, already installed in snapshot)
- Go >= 1.25, Node >= 22, `protoc` (protobuf-compiler via apt), `protoc-gen-go` (via `go install`)
- `GOPATH/bin` must be on `PATH` (added to `~/.bashrc`)

### Running the application
See `docs/commands.md` for full command reference. Key commands:

- **Dev mode (recommended):** `make devserver && ./wowsimtbc --usefs=true --launch=false --host=":3333"` serves the app at `http://localhost:3333/tbc/`
- **Full dev with Vite HMR + Go server:** `WATCH=1 make devmode` (or `npm start`) — starts Vite on port 5173 and Go server on port 3333. Requires `air` (`make setup` installs it).
- **Build dist first** if the `dist/tbc/` directory is empty: `make dist/tbc/.dirstamp`
- **Proto generation:** `make proto` — required after changing `.proto` files

### Testing
- **Go tests:** `make test` (builds WASM + runs `go test ./sim/...` with `--tags=with_db`)
- **TypeScript type-check:** `npm run type-check`
- **CSS lint:** `npm run lint:css`

### Gotchas
- The ESLint config (`eslint.config.mjs`) uses legacy format in a flat-config filename. `npm run lint:js` fails with a config error — this is a pre-existing repo issue.
- `make test` requires the WASM binary (`dist/tbc/lib.wasm`) and `binary_dist/dist.go`. The `test` target builds these automatically.
- The Go dev server (`wowsimtbc`) serves from `dist/` when `--usefs=true`. You must build dist content before it can serve pages.
- `protoc-gen-go` must be on `PATH` (lives in `$(go env GOPATH)/bin`). Ensure `export PATH=$PATH:$(go env GOPATH)/bin` is in your shell profile.
