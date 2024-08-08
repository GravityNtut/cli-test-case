package dataproductpurge

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
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

func PublishEventCommand() error {
	payload := `{"id":1,"uid":1,"name":"test","price":100,"kcal":50}`
	cmd := exec.Command(testutils.GravityCliString, "pub", "test", payload)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func CheckDataProductEventAmount(dataProduct string, amount string) error {
	cmd := testutils.GravityCliString + " product list"
	eventIndex := 3 // Because DESCRIPTION is null
	err := ut.ExecuteShell(cmd)
	if err != nil {
		return err
	}

	outStr := ut.CmdResult.Stdout

	outStrList := strings.Split(outStr, "\n")

	for _, outStr := range outStrList {
		if strings.Contains(outStr, dataProduct) {
			productItems := strings.Fields(outStr)

			if productItems[eventIndex] == amount {
				return nil
			}

			return errors.New("CLI returns error message: data product's event amount mismatches")
		}
	}

	return errors.New("CLI returns error message: data products doesn't contain the name " + dataProduct)
}

func CheckDataProductMessagesAmountByJetstream(dataProduct string, amount string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamsInfo()

	for stream := range streams {
		if stream.Config.Name == "GVT_default_DP_"+dataProduct {
			msgAmount, _ := strconv.ParseUint(amount, 10, 64)
			if stream.State.Msgs == msgAmount {
				return nil
			}

			return errors.New("CLI returns error message: stream's messages amount mismatches")
		}
	}
	return errors.New("CLI returns error message: streams doesn't contain the name GVT_default_DP_" + dataProduct)
}

func PurgeDataProduct(productName string) error {
	purgeString := "../gravity-cli product purge " + productName
	if err := ut.ExecuteShell(purgeString); err != nil {
		return err
	}
	return nil
}
func PurgeDataProductSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return errors.New("the data product purge should success")
}
func PurgeDataProductFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("the data product purge should fail")
}

func CheckError(errMsg string) error {
	outStr := ut.CmdResult.Stderr
	if strings.Contains(outStr, errMsg) {
		return nil
	}
	return errors.New("the error message is incorrect")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^Nats has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product with "'(.*?)'" using parameters "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^Create data product ruleset with "'(.*?)'", "'(.*?)'" using parameters "'(.*?)'"$`, ut.CreateDataProductRuleset)
	ctx.Given(`^Publish an Event$`, PublishEventCommand)
	ctx.Then(`^Check data product "'(.*?)'"'s Events amount is "'(.*?)'"$`, CheckDataProductEventAmount)
	ctx.Then(`^Use NATS JetStream to query the Messages amount of the data product "'(.*?)'" to be "'(.*?)'"$`, CheckDataProductMessagesAmountByJetstream)
	ctx.When(`^Purge data product "'(.*?)'"$`, PurgeDataProduct)
	ctx.Then(`^Check purging data product success$`, PurgeDataProductSuccess)
	ctx.Then(`^CLI returns exit code 1$`, PurgeDataProductFailed)
	ctx.Then(`^The error message should be "'(.*?)'"$`, CheckError)
}
