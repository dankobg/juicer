# Juicer (chess platform)

Project local setup:

1. `docker compose up -d`

2. `just mg-up`

3. `just keto-create-tuples`

4. `go run main.go identities import-identities`

5. `just certs-trust`

Run project:

1. `just dev`

2. `cd web && pnpm dev`
