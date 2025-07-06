// // #go:generate protoc -I=./ --go_out=paths=source_relative:./pb ./proto/juicer/juicer.proto
// // #go:generate go tool oapi-codegen github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config=api/oapi-codegen.yaml api/schema/juicer.yaml
// // #go:generate go tool jet -dsn={{jet_dsn}} -schema=public -path=./db/gen -ignore-tables=atlas_schema_revisions,continuity_containers,courier_message_dispatches,courier_messages,identities,identity_credential_identifiers,identity_credential_types,identity_credentials,identity_login_codes,identity_recovery_addresses,identity_recovery_codes,identity_recovery_tokens,identity_registration_codes,identity_verifiable_addresses,identity_verification_codes,identity_verification_tokens,keto_relation_tuples,keto_uuid_mappings,networks,schema_migration,selfservice_errors,selfservice_login_flows,selfservice_recovery_flows,selfservice_registration_flows,selfservice_settings_flows,selfservice_verification_flows,session_devices,session_token_exchanges,sessions

package main

import (
	"embed"
	"log"

	"github.com/dankobg/juicer/cmd"
)

//go:embed public/*
var publicFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func main() {
	if err := cmd.Run(publicFiles, templateFiles); err != nil {
		log.Fatalf("failed to run juicer chess server, %v", err)
	}
}
