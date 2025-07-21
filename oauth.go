package mbz

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tokenURLECE      = "https://ssoalpha.dvb.corpinter.net/v1/token"
	tokenURLAMAPNA   = "https://ssoalpha.am.dvb.corpinter.net/v1/token"
	audienceScopeECE = "audience:server:client_id:95B37AC2-D501-4CFD-B853-7D299DD2D872"
	audienceScopeAM  = "audience:server:client_id:87012BCA-0B2E-4127-BE24-97A71C1F3262"
)

func NewOAuth2Config(region Region, clientID, clientSecret string) (clientcredentials.Config, error) {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "groups", "profile"},
		TokenURL:     tokenURLECE,
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	switch region {
	case RegionECE:
		config.TokenURL = tokenURLECE
		config.Scopes = append(config.Scopes, audienceScopeECE)
	case RegionAMAPNA:
		config.TokenURL = tokenURLAMAPNA
		config.Scopes = append(config.Scopes, audienceScopeAM)
	default:
		return clientcredentials.Config{}, fmt.Errorf("invalid region: %s", region)
	}
	return config, nil
}
