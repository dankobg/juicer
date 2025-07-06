set dotenv-filename := ".env"
set dotenv-load

cwd := justfile_directory()

dev_mode := "nix" # nix or docker
domain := if env('DOMAIN', "") == "" { "juicer-dev.xyz" } else { '-a "${DOMAIN}"' }
redis_host := if env('REDIS_HOST', "") == "" { "localhost" } else { '-a "${REDIS_HOST}"' }
redis_port := if env('REDIS_PORT', "") == "" { "6379" } else { '-a "${REDIS_PORT}"' }
redis_pwd := if env('REDIS_PASSWORD', "") == "" { "" } else { '-a "${REDIS_PASSWORD}"' }
db_uri := if env('DB_URI', "") == "" { "postgres://test:test@localhost:5432/test?sslmode=disable" } else { '-a "${DB_URI}"' }
jet_dsn := if env('JET_DSN', "") == "" { db_uri } else { '-a "${JET_DSN}"' }
atlas_dsn := if env('ATLAS_DSN', "") == "" { db_uri + "&search_path=public" } else { '-a "${ATLAS_DSN}"' }
atlas_dev_dsn := if env('ATLAS_DEV_DSN', "") == "" { "postgres://test:test@localhost:5432/test_atlas?sslmode=disable&search_path=public" } else { '-a "${ATLAS_DEV_DSN}"' }
atlas_cmd := if dev_mode == "nix" { "atlas" } else { "docker compose run --rm atlas" }
migrations_dir := if dev_mode == "nix" { "file://db/migrations" } else { "file://migrations" }
migrations_schema := if dev_mode == "nix" { "file://db/schema.sql" } else { "file://schema.sql" }

migrations_format := (
	"""
	'{{ sql . "  " }}' 
	"""
)

sql_drop_public_tables := (
	"""
	DO $\\$
	DECLARE
		current_table text;
	BEGIN
		FOR current_table IN (SELECT table_name FROM information_schema.tables WHERE table_schema = 'public')
		LOOP
			EXECUTE 'DROP TABLE IF EXISTS public.' || current_table || ' CASCADE';
			END LOOP;
	END $\\$;
	"""
	)

# ----------------------------------------------------------------------------

default: 
	@just -l

# install needed tools
tools-install:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go get -tool github.com/go-jet/jet/v2/cmd/jet@latest

# generate jet db models, enums, tables etc.
gen-jet:
  jet -dsn={{jet_dsn}} -schema=public -path=./db/gen -ignore-tables=atlas_schema_revisions,continuity_containers,courier_message_dispatches,courier_messages,identities,identity_credential_identifiers,identity_credential_types,identity_credentials,identity_login_codes,identity_recovery_addresses,identity_recovery_codes,identity_recovery_tokens,identity_registration_codes,identity_verifiable_addresses,identity_verification_codes,identity_verification_tokens,keto_relation_tuples,keto_uuid_mappings,networks,schema_migration,selfservice_errors,selfservice_login_flows,selfservice_recovery_flows,selfservice_registration_flows,selfservice_settings_flows,selfservice_verification_flows,session_devices,session_token_exchanges,sessions

# Generate protobuf
gen-proto:
	protoc -I=./ --go_out=paths=source_relative:./pb ./proto/juicer/juicer.proto
	cd web && pnpm bufgen

# Generate openapi server
gen-openapi:
	oapi-codegen --config=api/oapi-codegen.yaml api/schema/juicer.yaml
	cd web && pnpm openapigen

# Run code generation
gen:
	go generate ./...
	just gen-proto
	just gen-openapi
	just gen-jet

# Run go vet
vet:
	go vet ./...

# Run go mod tidy
tidy:
	go mod tidy

# Show dependencies
deps:
	go mod graph

# Lint code with golangci-lint
lint:
	golangci-lint run

# Lint every dockerfile in docker dir
lint-dockerfiles:
	docker container run -v {{cwd}}/docker:/dockerfiles --rm -i hadolint/hadolint hadolint --ignore DL3018 /dockerfiles/app.dockerfile /dockerfiles/devspace.dockerfile /dockerfiles/goreleaser.dockerfile

# Run go fmt
fmt:
	go fmt ./...

# Run all tests
test:
	go test ./...

# Run all tests with verbose flag
test-verbose:
	go test -v ./...

# Run all tests with race flag
test-race:
	go test -race ./...

# Run all tests with race and verbose flags
test-racev:
	go test -v -race ./...

# Run test coverage
test-coverage:
	go test -cover ./...

# ----------------------------------------------------------------------------

# Run docker system prune (without unused volumes)
prune:
	docker system prune -a -f --volumes

# Run docker volume prune (annonymous and unused)
prune-volumes:
	docker volume prune -a -f 

# Run docker system prune all
prune-all: prune && prune-volumes

# ----------------------------------------------------------------------------

# Run docker compose commands
dc *flags:
	docker compose {{flags}}

# Run docker compose watch
dc-watch:
	docker compose watch

# Run docker compose watch without building & starting services
dc-watch-only:
	docker compose watch --no-up

# Run docker compose up -d
dc-up *flags:
	docker compose up -d {{flags}}

# Run docker compose down
dc-down:
	docker compose down

# Run docker compose start
dc-start:
	docker compose start

# Run docker compose stop
dc-stop:
	docker compose stop

# Run docker compose restart
dc-restart:
	docker compose restart

# ----------------------------------------------------------------------------

# Shell into a docker compose container by service name
sh name:
	docker compose exec -it {{name}} sh

# connect to redis via redis-cli
rediscli:
	docker compose exec -it redis redis-cli -h "{{redis_host}}" -p "{{redis_port}}" {{ if redis_pwd == "" { "" } else { redis_pwd } }}
	
# connect to postgres via psql
psql:
	docker compose exec -it pg psql "{{db_uri}}"

# ----------------------------------------------------------------------------

# migrations hash the directory
mg-hash: 
	{{atlas_cmd}} migrate hash --dir "{{migrations_dir}}"

# migrations create new
mg-new name *flags:
	{{atlas_cmd}} migrate new {{flags}} {{name}} --dir "{{migrations_dir}}"

# migrations diff
mg-diff name *flags:
	{{atlas_cmd}} migrate diff {{name}} --dir "{{migrations_dir}}" --to "{{migrations_schema}}" --dev-url "{{atlas_dev_dsn}}" --format {{migrations_format}} {{flags}}

# migrations apply
mg-apply *flags:
	{{atlas_cmd}} migrate apply --dir "{{migrations_dir}}" --url "{{atlas_dsn}}" {{flags}}

# migrations down
mg-down *flags:
	{{atlas_cmd}} migrate down --dir "{{migrations_dir}}" --url "{{atlas_dsn}}" --dev-url "{{atlas_dev_dsn}}" {{flags}}

# migrations status
mg-status *flags:
	{{atlas_cmd}} migrate status --dir "{{migrations_dir}}" --url "{{atlas_dsn}}" {{flags}}

# migrations validate
mg-validate *flags:
	{{atlas_cmd}} migrate validate --dir "{{migrations_dir}}" --dev-url "{{atlas_dev_dsn}}" {{flags}}

# schema command
mg-schema *flags:
	{{atlas_cmd}} schema --url "{{atlas_dsn}}" {{flags}}

# seed data
mg-seed:
	docker compose exec -T pg psql "{{db_uri}}" -f /seeds/seed_common.sql
	
# start from scratch in dev (clean all -> migrate -> seed_common)
mg-fresh:
	just mg-schema clean
	rm -rf ./db/migrations/*.sql
	just mg-hash
	just mg-diff initial
	just mg-apply --allow-dirty
	just dc-up kratos-migrate keto-migrate
	just mg-seed

# cleanup games for development
cleanup-games:
	docker compose exec -it redis redis-cli -h "{{redis_host}}" -p "{{redis_port}}" {{ if redis_pwd == "" { "" } else { redis_pwd } }} --raw "flushall"
	docker compose exec -T pg psql "{{db_uri}}" -c "truncate table game cascade;"

# ----------------------------------------------------------------------------

# backup postgres with pg_dumpall
pg-backup:
	docker compose exec -it pg pg_dumpall -c -U test > {{cwd}}/db/dev_dump.sql

# restore postgres backup
pg-restore: pg-dropall
	docker compose exec -T pg psql "{{db_uri}}" < {{cwd}}/db/dev_dump.sql

# drop all tables in public schema
pg-dropall:
	docker compose exec -T pg psql "{{db_uri}}" -c "{{sql_drop_public_tables}}"

# ----------------------------------------------------------------------------

# Import kratos identities
kratos-import-identities:
	docker compose exec kratos /bin/sh -c "cd /etc/config/kratos/imports && kratos import identity employees.json customers.json -e http://localhost:4434"

# Create keto relation tuples
keto-create-tuples:
	docker compose exec keto /bin/sh -c "keto relation-tuple create /etc/config/keto/relation-tuples -c /etc/config/keto/keto.yaml --format json --insecure-disable-transport-security"

# ----------------------------------------------------------------------------

# checks if mkcert is installed
_require_mkcert:
	#!/usr/bin/env sh
	command -v mkcert >/dev/null 2>&1 || { echo >&2 "mkcert is required, please install it to work with certs"; exit 1; }

# install mkcert
certs-install: _require_mkcert
	mkcert -install

# uninstall mkcert
certs-uninstall: _require_mkcert
	mkcert -uninstall && rm -rf "$(mkcert -CAROOT)"

# generate certs
certs: _require_mkcert
	rm -f {{cwd}}/certs/local*.pem && \
	mkcert -cert-file /tmp/local-cert.pem -key-file /tmp/local-key.pem "{{domain}}" "*.{{domain}}" localhost 127.0.0.1 ::1 && \
	cp /tmp/local-{key,cert}.pem {{cwd}}/certs && \
	rm -f /tmp/local-{cert,key}.pem

# ----------------------------------------------------------------------------

#   k:dev:
#     desc: Build and apply dev manifests
#     cmds:
#       - task: h:install:cert-manager
#       - kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.7.1/experimental-install.yaml
#       - kubectl -n cert-manager wait --for condition=established --timeout=5s crd/certificaterequests.cert-manager.io crd/certificates.cert-manager.io crd/challenges.acme.cert-manager.io crd/clusterissuers.cert-manager.io crd/issuers.cert-manager.io crd/orders.acme.cert-manager.io
#       - kubectl -n cert-manager wait --for=condition=ready --timeout=15s pod -l app=webhook -l app.kubernetes.io/name=webhook -l app.kubernetes.io/instance=cert-manager
#       - kubectl -n gateway-system wait --for=condition=ready --timeout=15s pod -l name=gateway-api-admission-server
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/common/overlays/dev | kubectl apply -f -
#       - task: h:install:postgresql
#       - task: h:install:redis
#       - task: h:install:traefik
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/kratos/overlays/dev | kubectl apply -f -
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/juicer/overlays/dev | kubectl apply -f -

#   k:prod:
#     desc: Build and apply prod manifests
#     cmds:
#       - kubectl apply -f k8s/initial/traefik-crd.yaml -f k8s/initial/cert-manager.yaml
#       - kubectl -n cert-manager wait --for condition=established --timeout=15s crd.apiextensions.k8s.io/certificaterequests.cert-manager.io crd.apiextensions.k8s.io/certificates.cert-manager.io crd.apiextensions.k8s.io/challenges.acme.cert-manager.io crd.apiextensions.k8s.io/clusterissuers.cert-manager.io crd.apiextensions.k8s.io/issuers.cert-manager.io crd.apiextensions.k8s.io/orders.acme.cert-manager.io crd.apiextensions.k8s.io/gatewayclasses.gateway.networking.k8s.io crd.apiextensions.k8s.io/gateways.gateway.networking.k8s.io crd.apiextensions.k8s.io/httproutes.gateway.networking.k8s.io
#       - kubectl -n cert-manager wait --for=condition=ready --timeout=15s pod -l app.kubernetes.io/name=webhook
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/initial | kubectl apply -f -
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/postgresql/overlays/prod | kubectl apply -f -
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/redis/overlays/prod | kubectl apply -f -
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/kratos/overlays/prod | kubectl apply -f -
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/juicer/overlays/prod | kubectl apply -f -

#   # ----------------------------------------------------------------------------

#   cluster:create:
#     desc: Create new k3d cluster
#     cmds:
#       - k3d cluster create dev -p "8081:80@loadbalancer" --k3s-arg="--disable=traefik@server:0"

#   cluster:delete:
#     desc: Create new k3d cluster
#     cmds:
#       - k3d cluster delete dev

#   # ----------------------------------------------------------------------------

#   crd:update-traefik:
#     desc: Update traefik CRDs manually with force conflicts
#     cmds:
#       - kubectl apply --server-side --force-conflicts -k https://github.com/traefik/traefik-helm-chart/traefik/crds/

#   # ----------------------------------------------------------------------------

#   h:repo:add:
#     desc: Helm add repositories
#     cmds:
#       - helm repo add jetstack https://charts.jetstack.io
#       - helm repo add cnpg https://cloudnative-pg.github.io/charts
#       - helm repo add traefik https://traefik.github.io/charts
#       - helm repo add ory https://k8s.ory.sh/helm/charts
#       - helm repo add bitnami https://charts.bitnami.com/bitnami

#   h:repo:remove:
#     desc: Helm remove repositories
#     cmds:
#       - helm repo remove jetstack https://charts.jetstack.io
#       - helm repo remove cnpg https://cloudnative-pg.github.io/charts
#       - helm repo remove traefik https://traefik.github.io/charts
#       - helm repo remove ory https://k8s.ory.sh/helm/charts
#       - helm repo remove bitnami https://charts.bitnami.com/bitnami

#   h:repo:update:
#     desc: Helm update repositories
#     cmds:
#       - helm repo update

#   h:install:
#     desc: Helm install
#     cmds:
#       - task: h:install:cert-manager
#       - task: h:install:postgresql
#       - task: h:install:redis
#       - task: h:install:traefik

#   h:install:postgresql:
#     desc: Helm install postgresql
#     cmds:
#       - helm secrets install postgresql bitnami/postgresql -n dev --create-namespace -f k8s/postgresql/dev-values.yaml -f k8s/postgresql/dev-secrets.yaml

#   h:install:redis:
#     desc: Helm install redis
#     cmds:
#       - helm secrets install redis bitnami/redis -n dev --create-namespace -f k8s/redis/dev-values.yaml -f k8s/redis/dev-secrets.yaml

#   h:install:cert-manager:
#     desc: Helm install cert manager
#     cmds:
#       - helm install cert-manager jetstack/cert-manager -n cert-manager --create-namespace -f k8s/cert-manager/dev-values.yaml

#   h:install:traefik:
#     desc: Helm install traefik
#     cmds:
#       # - helm install traefik traefik/traefik -n traefik --create-namespace -f k8s/traefik/dev-values.yaml
#       - helm install traefik traefik/traefik -n dev --create-namespace -f k8s/traefik/dev-values.yaml

#   h:delete:
#     desc: Helm delete
#     cmds:
#       - task: h:delete:cert-manager
#       - task: h:delete:traefik
#       - task: h:delete:postgresql
#       - task: h:delete:redis

#   h:delete:postgresql:
#     desc: Helm delete postgresql
#     cmds:
#       - helm delete postgresql -n dev

#   h:delete:redis:
#     desc: Helm delete redis
#     cmds:
#       - helm delete redis -n dev

#   h:delete:cert-manager:
#     desc: Helm delete cert manager
#     cmds:
#       - helm delete cert-manager -n cert-manager

#   h:delete:traefik:
#     desc: Helm delete traefik
#     cmds:
#       - helm delete traefik -n traefik

#   # ----------------------------------------------------------------------------

#   k:wat:
#     cmds:
#       - task: h:install:cert-manager
#       - kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.7.1/experimental-install.yaml
#       - kubectl -n cert-manager wait --for condition=established --timeout=5s crd/certificaterequests.cert-manager.io crd/certificates.cert-manager.io crd/challenges.acme.cert-manager.io crd/clusterissuers.cert-manager.io crd/issuers.cert-manager.io crd/orders.acme.cert-manager.io
#       - kubectl -n cert-manager wait --for=condition=ready --timeout=15s pod -l app=webhook -l app.kubernetes.io/name=webhook -l app.kubernetes.io/instance=cert-manager
#       - kubectl -n gateway-system wait --for=condition=ready --timeout=15s pod -l name=gateway-api-admission-server
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/common/overlays/dev | kubectl apply -f -
#       - task: h:install:traefik
#       - kustomize build --enable-alpha-plugins --enable-exec k8s/juicer/overlays/dev | kubectl apply -f -
