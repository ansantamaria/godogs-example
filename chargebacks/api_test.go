package main

import (
	"fmt"

	"github.com/cucumber/godog"
)

type apiFeature struct {
	status int
	err    error
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.status = 0
	a.err = nil
}

func (a *apiFeature) iSendRequestTo(endpoint string) (err error) {
	switch endpoint {
	case "resolveChargeback":
		a.status, a.err = resolveChargeback()
		if a.err != nil {
			err = fmt.Errorf("error: %s", a.err)
		}
		fmt.Sprint(a.status)
	default:
		err = fmt.Errorf("unknown endpoint: %s", endpoint)
	}
	return
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.status {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.status)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.BeforeScenario(api.resetResponse)

	ctx.Step(`^i send PUT request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)

}
