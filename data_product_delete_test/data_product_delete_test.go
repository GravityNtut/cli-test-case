package dataproductdelete

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

func TestMain(m *testing.M) {
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

func SearchDataProductByCLIFail(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return errors.New("the CLI was expected to not exist, but it was found")
}

func SearchDataProductByJetstreamFail(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return errors.New("jetstream was expected to not exist, but it was found")
		}
	}
	return nil
}

func DeleteDataProductCommand(dataProduct string) error {
	commandString := "../gravity-cli product delete "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDataProductSuccess(productName string) error {
	outStr := ut.CmdResult.Stdout
	if outStr == "Product \""+productName+"\" was deleted\n" {
		return nil
	}
	return errors.New("the message returned by the CLI is incorrect")
}

func DeleteDataProductFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("the data product deletion should fail")
}

// TODO
// func AssertErrorMessages(errorMessage string) error {
// 	outErr := ut.CmdResult.Stderr
// 	if outErr == errorMessage {
// 		return nil
// 	}
// 	return errors.New("CLI returns error message")
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.RestartDocker()
		return ctx, nil
	})

	ctx.Given(`^Nats has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product with "'(.*?)'" using parameters "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.When(`^Delete data product "'(.*?)'"$`, DeleteDataProductCommand)
	ctx.Then(`^The CLI returned the message: Product "'(.*?)'" was deleted.$`, DeleteDataProductSuccess)
	ctx.Then(`^CLI returns exit code 1$`, DeleteDataProductFail)
	ctx.Then(`^Using gravity-cli to query "'(.*?)'" shows it does not exist.$`, SearchDataProductByCLIFail)
	ctx.Then(`^Using NATS Jetstream to query "'(.*?)'" shows it does not exist.$`, SearchDataProductByJetstreamFail)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, AssertErrorMessages)
}
