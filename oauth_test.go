package mbz

import (
	"testing"

	"golang.org/x/oauth2"
)

func TestNewOAuth2Config(t *testing.T) {
	const clientID = "test-client-id"
	const clientSecret = "test-client-secret"
	tests := []struct {
		name             string
		region           Region
		expectedTokenURL string
		expectedScopes   []string
		expectError      bool
	}{
		{
			name:             "region ECE",
			region:           RegionECE,
			expectedTokenURL: TokenURLECE,
			expectedScopes:   []string{"openid", "groups", "profile", AudienceScopeECE},
			expectError:      false,
		},
		{
			name:             "region AMAP/NA",
			region:           RegionAMAPNA,
			expectedTokenURL: TokenURLAMAPNA,
			expectedScopes:   append([]string{"openid", "groups", "profile"}, AudienceScopeAMAPNA),
			expectError:      false,
		},
		{
			name:        "region empty",
			region:      "",
			expectError: true,
		},
		{
			name:        "region invalid",
			region:      "INVALID",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewOAuth2Config(tt.region, clientID, clientSecret)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for region %q, but got none", tt.region)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for region %q: %v", tt.region, err)
				return
			}
			if config.ClientID != clientID {
				t.Errorf("Expected ClientID %q, got %q", clientID, config.ClientID)
			}
			if config.ClientSecret != clientSecret {
				t.Errorf("Expected ClientSecret %q, got %q", clientSecret, config.ClientSecret)
			}
			if config.TokenURL != tt.expectedTokenURL {
				t.Errorf("Expected TokenURL %q, got %q", tt.expectedTokenURL, config.TokenURL)
			}
			if config.AuthStyle != oauth2.AuthStyleInParams {
				t.Errorf("Expected AuthStyle %v, got %v", oauth2.AuthStyleInParams, config.AuthStyle)
			}
			if len(config.Scopes) != len(tt.expectedScopes) {
				t.Errorf("Expected %d scopes, got %d", len(tt.expectedScopes), len(config.Scopes))
			}
			for i, expectedScope := range tt.expectedScopes {
				if i >= len(config.Scopes) || config.Scopes[i] != expectedScope {
					t.Errorf("Expected scope[%d] %q, got %q", i, expectedScope, config.Scopes[i])
				}
			}
		})
	}
}
