package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog/colors"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func thereAreGodogs(available int) error {
	Godogs = available
	return nil
}

func iEat(num int) error {
	err := assertExpectedAndActual(
		assert.GreaterOrEqual, Godogs, num,
		"You cannot eat %d godogs, there are %d available", num, Godogs,
	)
	if err != nil {
		return err
	}

	Godogs -= num
	return nil
}

func thereShouldBeRemaining(remaining int) error {
	return assertExpectedAndActual(
		assert.Equal, Godogs, remaining,
		"Expected %d godogs to be remaining, but there is %d", remaining, Godogs,
	)
}

func thereShouldBeNoneRemaining() error {
	return assertActual(
		assert.Empty, Godogs,
		"Expected none godogs to be remaining, but there is %d", Godogs,
	)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() { Godogs = 0 })
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		Godogs = 0
	})

	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.Step(`^I eat (\d+)$`, iEat)
	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
	ctx.Step(`^there should be none remaining$`, thereShouldBeNoneRemaining)
}

func assertExpectedAndActual(expectedAndActualAssertion expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	expectedAndActualAssertion(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

func assertActual(a actualAssertion, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, actual, msgAndArgs...)
	return t.err
}

type actualAssertion func(t assert.TestingT, actual interface{}, msgAndArgs ...interface{}) bool

type asserter struct {
	err error
}

func (asserter *asserter) Errorf(format string, args ...interface{}) {
	asserter.err = fmt.Errorf(format, args...)
}
