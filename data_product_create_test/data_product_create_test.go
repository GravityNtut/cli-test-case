package dataproductcreate

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

var ut = testutils.TestUtils{}

var opts = godog.Options{
	Format:        "pretty",
	Paths:         []string{"./"},
	StopOnFailure: ut.Config.StopOnFailure,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(_ *testing.M) {
	pflag.Parse()
	err := ut.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}
	if suite.Run() != 0 {
		log.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func CreateDataProductCommand(dataProduct string, description string, schema string, enabled string) error {
	dataProduct = ut.ProcessString(dataProduct)
	commandString := "../gravity-cli product create "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	if description != testutils.IgnoreString {
		description := ut.ProcessString(description)
		commandString += " --desc " + description
	}

	if schema != testutils.IgnoreString {
		commandString += " --schema " + schema
	}

	if enabled == testutils.TrueString {
		commandString += " --enabled"
	} else if enabled == testutils.FalseString {
		commandString += " --enabled=false"
	} else if enabled != testutils.IgnoreString {
		return errors.New("invalid parameter for enable")
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func CreateDataProductCommandSuccess(productName string) error {
	outStr := ut.CmdResult.Stdout
	productName = ut.ProcessString(productName)
	if outStr == "Product \""+productName+"\" was created\n" {
		return nil
	}
	return errors.New("the message returned by the CLI is incorrect")
}

func CreateDataProductCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("the message returned by the CLI is incorrect")
}

func SearchDataProductByCLISuccess(dataProduct string) error {
	dataProduct = ut.ProcessString(dataProduct)
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	return err
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	dataProduct = ut.ProcessString(dataProduct)
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return nil
		}
	}
	return errors.New("data product creation in Jetstream was unsuccessful")
}

// TODO
// func AssertErrorMessages(errorMessage string) error {
// 	return nil
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^NATS has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.When(`^Create data product "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, CreateDataProductCommand)
	ctx.Then(`^Cli returns "'(.*?)'" created successfully$`, CreateDataProductCommandSuccess)
	ctx.Then(`^CLI returns exit code 1$`, CreateDataProductCommandFail)
	ctx.Then(`^Using gravity-cli to query "'(.*?)'" shows it exist$`, SearchDataProductByCLISuccess)
	ctx.Then(`^Using NATS Jetstream to query "'(.*?)'" shows it exist$`, SearchDataProductByJetstreamSuccess)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, AssertErrorMessages)
}
