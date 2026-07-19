package shared

import "fmt"

func AuthzIdentityID(id string) string {
	return fmt.Sprintf("identity:%s", id)
}

func AuthzCourierMessageID(id string) string {
	return fmt.Sprintf("courier_message:%s", id)
}

func AuthzSchemaID(id string) string {
	return fmt.Sprintf("schema:%s", id)
}

func AuthzSessionID(id string) string {
	return fmt.Sprintf("session:%s", id)
}
