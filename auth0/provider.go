package auth0

import (
	"context"
	"fmt"
	"os"

	"github.com/alekc/terraform-provider-auth0/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	"gopkg.in/auth0.v5"
	"gopkg.in/auth0.v5/management"
)

var provider *schema.Provider

func init() {
	provider = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_DOMAIN", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTH0_CLIENT_SECRET", nil),
			},
			"debug": {
				Type:     schema.TypeBool,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					v := os.Getenv("AUTH0_DEBUG")
					if v == "" {
						return false, nil
					}
					return v == "1" || v == "true" || v == "on", nil
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"auth0_client":          newClient(),
			"auth0_client_grant":    newClientGrant(),
			"auth0_connection":      newConnection(),
			"auth0_custom_domain":   newCustomDomain(),
			"auth0_resource_server": newResourceServer(),
			"auth0_rule":            newRule(),
			"auth0_rule_config":     newRuleConfig(),
			"auth0_hook":            newHook(),
			"auth0_prompt":          newPrompt(),
			"auth0_email":           newEmail(),
			"auth0_email_template":  newEmailTemplate(),
			"auth0_user":            newUser(),
			"auth0_tenant":          newTenant(),
			"auth0_role":            newRole(),
			"auth0_log_stream":      newLogStream(),
			"auth0_branding":        newBranding(),
			"auth0_guardian":        newGuardian(),
			"auth0_action":          newAction(),
			"auth0_flow":            newFlow(),
		},
		ConfigureContextFunc: Configure,
	}
}

func Provider() *schema.Provider {
	return provider
}

func Configure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	domain := data.Get("domain").(string)
	id := data.Get("client_id").(string)
	secret := data.Get("client_secret").(string)
	debug := data.Get("debug").(bool)

	userAgent := fmt.Sprintf("Terraform-Provider-Auth0/%s (Go-Auth0-SDK/%s; Terraform-SDK/%s; Terraform/%s)",
		Version(),
		SDKVersion(),
		TerraformSDKVersion(),
		TerraformVersion())

	m, err := management.New(domain,
		management.WithClientCredentials(id, secret),
		management.WithDebug(debug),
		management.WithUserAgent(userAgent))
	return m, diag.FromErr(err)
}

func Version() string {
	return version.ProviderVersion
}

func SDKVersion() string {
	return auth0.Version
}

func TerraformVersion() string {
	return provider.TerraformVersion
}

func TerraformSDKVersion() string {
	return meta.SDKVersionString()
}
