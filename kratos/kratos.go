package kratos

import kratos "github.com/ory/client-go"

type Client struct {
	Public *kratos.APIClient
	Admin  *kratos.APIClient
}

func NewClient(publicURL, adminURL string) *Client {
	publicConf := kratos.NewConfiguration()
	adminConf := kratos.NewConfiguration()

	publicConf.Servers = kratos.ServerConfigurations{{
		URL:         publicURL,
		Description: "Kratos public server",
		Variables: map[string]kratos.ServerVariable{
			"api": {
				Description:  "Target the public API.",
				DefaultValue: "public",
				EnumValues:   []string{"public"},
			},
			"tenant": {
				Description:  "Juicer tenant",
				DefaultValue: "juicer",
			},
		},
	}}

	adminConf.Servers = kratos.ServerConfigurations{{
		URL:         adminURL,
		Description: "Kratos admin server",
		Variables: map[string]kratos.ServerVariable{
			"api": {
				Description:  "Target the administrative API.",
				DefaultValue: "admin",
				EnumValues:   []string{"admin"},
			},
			"tenant": {
				Description:  "Juicer tenant",
				DefaultValue: "juicer",
			},
		},
	}}

	public := kratos.NewAPIClient(publicConf)
	admin := kratos.NewAPIClient(adminConf)

	c := &Client{
		Public: public,
		Admin:  admin,
	}

	return c
}
