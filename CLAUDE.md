# Wallet

Crypto wallet dashboard for monitoring EVM RPC endpoints and executing JSON-RPC calls.

## Project Structure

- `cmd/wallet/` — Entry point
- `internal/config/` — Environment config
- `internal/endpoint/` — Endpoint store (JSON file), RPC polling, CRUD
- `internal/server/` — Echo HTTP server, routes, dashboard

## Build & Run

```bash
go build -o wallet ./cmd/wallet
go vet ./...

# Docker
./.launch.sh
```

## Conventions

- Go module: `github.com/primal-host/wallet`
- HTTP framework: Echo v4
- Container name: `crypto-wallet`
- No database — endpoints stored in `endpoints.json` file
- Config uses env vars (`LISTEN_ADDR`, `ENDPOINTS_FILE`)

## Docker

- Image/container: `crypto-wallet`
- Network: `infra` (traefik)
- Port: 4322
- Traefik: `wallet.primal.host` / `wallet.localhost`
- Traefik middleware: `noknok-auth@docker` (AT Protocol OAuth via noknok)
- DNS: `192.168.147.53` (infra CoreDNS)
- `endpoints.json` mounted as volume and baked into image at `/etc/wallet/endpoints.json`

## Authentication

All access is gated by noknok forwardAuth via Traefik. No internal auth — the app trusts that Traefik only forwards authenticated requests.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/` | Dashboard |
| `GET` | `/api/status` | Poll all endpoints (chain ID, block number, latency) |
| `POST` | `/api/rpc/:id` | Proxy JSON-RPC call to named endpoint |
| `POST` | `/api/endpoints` | Add endpoint (name, url, symbol) |
| `PUT` | `/api/endpoints/:id` | Update endpoint |
| `DELETE` | `/api/endpoints/:id` | Delete endpoint |

## Endpoint Store

Endpoints are loaded from `endpoints.json` at startup. CRUD operations persist back to the same file. Each endpoint has:

- `id` — auto-generated slug from name
- `name` — display name (e.g., "Avalanche C-Chain")
- `url` — RPC URL (may include basic auth credentials)
- `symbol` — native token symbol (e.g., "AVAX", "ETH")

Polling calls `eth_chainId` and `eth_blockNumber` to check liveness and latency.
