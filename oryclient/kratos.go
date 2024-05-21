package oryclient

// import kratos "github.com/ory/client-go"

// type KratosClient struct {
// 	PublicApi *kratos.APIClient
// 	AdminApi  *kratos.APIClient
// }

// func NewKratosClient(publicURL, adminURL string) *KratosClient {
// 	publicConf := kratos.NewConfiguration()
// 	adminConf := kratos.NewConfiguration()

// 	publicConf.Servers = kratos.ServerConfigurations{{
// 		URL:         publicURL,
// 		Description: "Kratos public server",
// 		Variables: map[string]kratos.ServerVariable{
// 			"api": {
// 				Description:  "Target the public API.",
// 				DefaultValue: "public",
// 				EnumValues:   []string{"public"},
// 			},
// 			"tenant": {
// 				Description:  "Tenant ID as provided by Ory Cloud.",
// 				DefaultValue: "animond",
// 			},
// 		},
// 	}}

// 	adminConf.Servers = kratos.ServerConfigurations{{
// 		URL:         adminURL,
// 		Description: "Kratos admin server",
// 		Variables: map[string]kratos.ServerVariable{
// 			"api": {
// 				Description:  "Target the administrative API.",
// 				DefaultValue: "admin",
// 				EnumValues:   []string{"admin"},
// 			},
// 			"tenant": {
// 				Description:  "Tenant ID as provided by Ory Cloud.",
// 				DefaultValue: "animond",
// 			},
// 		},
// 	}}

// 	public := kratos.NewAPIClient(publicConf)
// 	admin := kratos.NewAPIClient(adminConf)

// 	c := &KratosClient{
// 		PublicApi: public,
// 		AdminApi:  admin,
// 	}

// 	return c
// }
