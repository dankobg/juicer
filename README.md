# Juicer (chess platform)

Run caddy devproxy

1. `dc -p devproxy -f compose.devproxy.yaml up -d`

Juicer local setup:

1. `docker compose up -d`

2. `just mg-up`

3. `just keto-create-tuples`

4. `go run main.go identities import-identities`

5. `just devproxy-setup`

6. `just certs-trust`


Run juicer:

1. `just dev`

2. `cd web && pnpm dev`
