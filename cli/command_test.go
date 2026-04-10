package cli

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/way-platform/mbz-go"
	"golang.org/x/oauth2"
)

func TestResolveOAuth2RegionUsesStoredRegionFirst(t *testing.T) {
	t.Parallel()

	creds := &FleetCredentials{
		Region: string(mbz.RegionAMAPNA),
	}
	region, err := resolveOAuth2Region(
		creds,
		&oauth2.Token{AccessToken: testJWT("https://ssoalpha.dvb.corpinter.net/v1")},
	)
	if err != nil {
		t.Fatalf("resolve region: %v", err)
	}
	if region != mbz.RegionAMAPNA {
		t.Fatalf("got region %q, want %q", region, mbz.RegionAMAPNA)
	}
}

func TestResolveOAuth2RegionInfersFromTokenIssuer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		iss  string
		want mbz.Region
	}{
		{name: "ece", iss: "https://ssoalpha.dvb.corpinter.net/v1", want: mbz.RegionECE},
		{name: "amapna", iss: "https://ssoalpha.am.dvb.corpinter.net/v1", want: mbz.RegionAMAPNA},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			region, err := resolveOAuth2Region(&FleetCredentials{}, &oauth2.Token{
				AccessToken: testJWT(tt.iss),
			})
			if err != nil {
				t.Fatalf("resolve region: %v", err)
			}
			if region != tt.want {
				t.Fatalf("got region %q, want %q", region, tt.want)
			}
		})
	}
}

func TestInferRegionFromAccessTokenErrorsOnUnknownIssuer(t *testing.T) {
	t.Parallel()

	_, err := inferRegionFromAccessToken(testJWT("https://example.com"))
	if err == nil {
		t.Fatal("expected error for unknown issuer")
	}
}

func testJWT(issuer string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"iss":%q}`, issuer)))
	return header + "." + payload + "."
}
