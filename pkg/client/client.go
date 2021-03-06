// Package client provides intergation with api calls.
package client

import (
	"net/http"
	"time"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps/app_rules"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers/access_token_claims"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers/scopes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/legal_values"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/roles"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/session_login_tokens"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/user_mappings"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"
)

// APIClient is used to communicate with the available api services.
type APIClient struct {
	clientID     string
	clientSecret string
	region       string
	baseURL      string
	client       *http.Client
	Services     *Services
}

// Services contains all the available api services.
type Services struct {
	HTTPService          *olhttp.OLHTTPService
	AppsV2               *apps.V2Service
	AppRulesV2           *apprules.V2Service
	UsersV2              *users.V2Service
	UserMappingsV2       *usermappings.V2Service
	SessionLoginTokensV1 *sessionlogintokens.V1Service
	AuthServersV2        *authservers.V2Service
	AccessTokenClaimsV2  *accesstokenclaims.V2Service
	ScopesV2             *scopes.V2Service
	SmartHooksV1         *smarthooks.V1Service
	RolesV1              *roles.V1Service
}

// NewClient uses the config to generate the api client with services attached, and returns
// the new api client.
func NewClient(cfg *APIClientConfig) (*APIClient, error) {
	cfg, err := cfg.Initialize()
	if err != nil {
		return &APIClient{}, err
	}

	httpClient := &http.Client{
		Timeout: time.Second * time.Duration(cfg.Timeout),
	}

	resourceRepository := olhttp.New(services.HTTPServiceConfig{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		BaseURL:      cfg.Url,
		Client:       httpClient,
	})

	legalValuesService := legalvalues.New(resourceRepository, cfg.Url)

	return &APIClient{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		region:       cfg.Region,
		baseURL:      cfg.Url,
		client:       httpClient,
		Services: &Services{
			HTTPService:          resourceRepository,
			AppsV2:               apps.New(resourceRepository, cfg.Url),
			AppRulesV2:           apprules.New(resourceRepository, legalValuesService, cfg.Url),
			UserMappingsV2:       usermappings.New(resourceRepository, legalValuesService, cfg.Url),
			UsersV2:              users.New(resourceRepository, cfg.Url),
			SessionLoginTokensV1: sessionlogintokens.New(resourceRepository, cfg.Url),
			AuthServersV2:        authservers.New(resourceRepository, cfg.Url),
			AccessTokenClaimsV2:  accesstokenclaims.New(resourceRepository, cfg.Url),
			ScopesV2:             scopes.New(resourceRepository, cfg.Url),
			SmartHooksV1:         smarthooks.New(resourceRepository, cfg.Url),
			RolesV1:              roles.New(resourceRepository, cfg.Url),
		},
	}, nil
}
