package mbz

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// NewOAuth2Config creates a new OAuth2 [clientcredentials.Config] for the given region.
func NewOAuth2Config(region Region, clientID, clientSecret string) (clientcredentials.Config, error) {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "groups", "profile"},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	switch region {
	case RegionECE:
		config.TokenURL = TokenURLECE
		config.Scopes = append(config.Scopes, AudienceScopeECE)
	case RegionAMAPNA:
		config.TokenURL = TokenURLAMAPNA
		config.Scopes = append(config.Scopes, AudienceScopeAMAPNA)
	default:
		return clientcredentials.Config{}, fmt.Errorf("invalid region: %s", region)
	}
	return config, nil
}
