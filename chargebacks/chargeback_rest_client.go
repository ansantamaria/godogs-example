package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	commonContext "github.com/mercadolibre/fury_credits-go-commons/v3/context"
	"github.com/mercadolibre/fury_go-core/pkg/telemetry/tracing"
	"github.com/mercadolibre/fury_go-core/pkg/transport/httpclient"
)

const (
	chargebackResolveURI = "%s/sites/%s/purchases/%s/chargebacks/%s"
)

// ChargebackRestClient struct that represents the Rest client that has the purpose of going to the Chargeback.
type ChargebackRestClient struct {
	restClient httpclient.Requester
	urlBase    string
}

// NewChargebackRestClient constructor that allows from a configuration of the RequestBuilder to create a client to the
// Chargeback API.
func NewChargebackRestClient(restClient httpclient.Requester, urlBase string) *ChargebackRestClient {
	return &ChargebackRestClient{
		restClient: restClient,
		urlBase:    urlBase,
	}
}

// ResolveChargeback function that requests to chargeback api resolve a chargeback.
func (client ChargebackRestClient) ResolveChargeback(ctx context.Context, siteID string, purchaseID string) (int, error) {
	const trackingAlias = "CUCUMBER_TEST"
	logger := commonContext.Logger(ctx)
	url := fmt.Sprintf(chargebackResolveURI, client.urlBase, siteID, purchaseID, purchaseID)
	body := struct {
		Answer string `json:"answer"`
	}{Answer: "invalid"}
	bodyBytes, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(bodyBytes)

	ctx = tracing.WithTargetID(ctx, trackingAlias)
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bodyReader)
	request.Header.Set("X-Client-Id", "2802671000854551")

	if err != nil {
		logger.Errorf(err.Error())
		return http.StatusInternalServerError, err
	}

	response, err := client.restClient.Do(request)

	if err != nil {
		logger.Errorf(err.Error())
		return http.StatusInternalServerError, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Errorf("status code :%s", response.StatusCode)
		return response.StatusCode, nil
	}
	responseDecode := struct {
		Message string `json:"message"`
	}{Message: ""}
	err = json.NewDecoder(response.Body).Decode(&responseDecode)

	if err != nil {
		logger.Errorf(err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
