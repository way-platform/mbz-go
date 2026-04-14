package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/way-platform/mbz-go/api/servicesv1"
	fleetv1 "github.com/way-platform/mbz-go/proto/gen/go/wayplatform/connect/mercedesbenz/fleet/v1"
)

// ListServices lists the available data services for the account.
func (c *Client) ListServices(
	ctx context.Context,
	request *fleetv1.ListServicesRequest,
) (_ *fleetv1.ListServicesResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("mbz: list services: %w", err)
		}
	}()
	path := "/v2/accounts/services"
	if request.GetDetails() {
		path = "/v2/accounts/services/details"
	}
	requestURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("User-Agent", getUserAgent())
	httpResponse, err := c.httpClient(c.config).Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := httpResponse.Body.Close(); closeErr != nil {
			log.Printf("mbz: failed to close response body: %v", closeErr)
		}
	}()
	if httpResponse.StatusCode != http.StatusOK {
		return nil, newResponseError(httpResponse)
	}
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	var responseBody servicesv1.GetAllServicesWithSignalsAndCommandsResponse
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return nil, err
	}
	services := make([]*fleetv1.Service, 0, len(responseBody.Services))
	for _, s := range responseBody.Services {
		ps := &fleetv1.Service{}
		ps.SetId(s.ID)
		ps.SetName(s.Name)
		ps.SetCiamScope(s.CIAMScope)
		ps.SetConsent(string(s.Consent))
		if s.CountryCode != "" {
			ps.SetCountryCode(s.CountryCode)
		}
		roles := make([]string, 0, len(s.Roles))
		for _, r := range s.Roles {
			roles = append(roles, string(r))
		}
		ps.SetRoles(roles)
		signals := make([]*fleetv1.ServiceSignal, 0, len(s.Signals))
		for _, sig := range s.Signals {
			psig := &fleetv1.ServiceSignal{}
			psig.SetName(sig.Name)
			psig.SetDataType(string(sig.DataType))
			psig.SetMandatory(sig.Mandatory)
			if sig.Unit != "" {
				psig.SetUnit(string(sig.Unit))
			}
			behaviours := make([]string, 0, len(sig.SendingBehaviour))
			for _, b := range sig.SendingBehaviour {
				behaviours = append(behaviours, string(b))
			}
			psig.SetSendingBehaviour(behaviours)
			signals = append(signals, psig)
		}
		ps.SetSignals(signals)
		commands := make([]*fleetv1.ServiceCommand, 0, len(s.Commands))
		for _, cmd := range s.Commands {
			pcmd := &fleetv1.ServiceCommand{}
			pcmd.SetName(cmd.Name)
			pcmd.SetMandatory(cmd.Mandatory)
			commands = append(commands, pcmd)
		}
		ps.SetCommands(commands)
		services = append(services, ps)
	}
	resp := &fleetv1.ListServicesResponse{}
	resp.SetServices(services)
	if responseBody.Version != "" {
		resp.SetVersion(responseBody.Version)
	}
	return resp, nil
}
