package main

import (
	"context"
	"time"

	"github.com/mercadolibre/fury_go-core/pkg/transport"
	"github.com/mercadolibre/fury_go-core/pkg/transport/httpclient"
)

const (
	maxRetriesAccountRestClient     = 3
	connectTimeOutAccountRestClient = 300
	dialTimeOutRestClient           = 700
	ExponentialBackoffMin           = 10
	ExponentialBackoffMax           = 100
)

var rc *ChargebackRestClient

func resolveChargeback() (int, error) {
	rc = NewChargebackRestClient(ConfigurationChargebackRestClient(), "https://internal-api.mercadopago.com/credit-card-beta/")
	return rc.ResolveChargeback(context.TODO(), "MLB", "15a8785f-5c0a-11eb-9acb-0242ac120009")
}

func ConfigurationChargebackRestClient() *httpclient.RetryableClient {
	return getHTTPClient(connectTimeOutAccountRestClient, maxRetriesAccountRestClient)
}

func getHTTPClient(timeOut, maxRetries int) *httpclient.RetryableClient {
	tcpTimeout := time.Duration(dialTimeOutRestClient) * time.Millisecond
	clientTransport := transport.NewTransport(transport.OptionDialTimeout(tcpTimeout))
	clientPooledTransport := transport.NewPooledFromTransport("namet", clientTransport)
	httpClientTimeout := time.Duration(timeOut) * time.Millisecond

	opts := []httpclient.OptionRetryable{
		httpclient.WithTransport(clientPooledTransport),
		httpclient.WithTimeout(httpClientTimeout),
		httpclient.WithBackoffStrategy(
			httpclient.ExponentialBackoff(
				time.Duration(ExponentialBackoffMin)*time.Millisecond,
				time.Duration(ExponentialBackoffMax)*time.Millisecond),
		),
	}
	return httpclient.NewRetryable(maxRetries, opts...)
}
