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
)

var ut = testutils.TestUtils{}

func TestFeatures(t *testing.T) {
	err := ut.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"./"},
			StopOnFailure: ut.Config.StopOnFailure,
			TestingT:      t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func AddRulesetCommand(dataProduct string, ruleset string, event string, enabled string) error {
	if enabled == testutils.TrueString {
		cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--enabled", "--event", event, "--method", "create", "-s", ut.Config.JetstreamURL)
		return cmd.Run()
	} else if enabled == testutils.IgnoreString || enabled == testutils.FalseString {
		cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", event, "--method", "create", "-s", ut.Config.JetstreamURL)
		return cmd.Run()
	}
	return errors.New("the parameters of the ruleset add are incorrect")
}

func PublishEventCommand(event string, payload string) error {
	pubString := "../gravity-cli pub " + event + " " + payload
	if err := ut.ExecuteShell(pubString); err != nil {
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
	commandString := "../gravity-cli product purge"
	if productName != testutils.NullString && productName != testutils.IgnoreString {
		commandString += " " + productName
	}

	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func PurgeDataProductFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("the data product purge should fail")
}

// func CheckError(errMsg string) error {
// 	outStr := ut.CmdResult.Stderr
// 	if strings.Contains(outStr, errMsg) {
// 		return nil
// 	}
// 	return errors.New(errMsg)
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.RestartDocker()
		return ctx, nil
	})

	ctx.Given(`^Nats has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product with "'(.*?)'" using parameters "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^Create data product ruleset with "'(.*?)'", "'(.*?)'" using parameters "'(.*?)'", "'(.*?)'"$`, AddRulesetCommand)
	ctx.Given(`^Publish Event "'(.*?)'" using parameters "'(.*?)'"$`, PublishEventCommand)
	ctx.Then(`^Check data product "'(.*?)'"'s Events amount is "'(.*?)'"$`, CheckDataProductEventAmount)
	ctx.Then(`^Use NATS JetStream to query the Messages amount of the data product "'(.*?)'" to be "'(.*?)'"$`, CheckDataProductMessagesAmountByJetstream)
	ctx.When(`^Purge data product "'(.*?)'"$`, PurgeDataProduct)
	ctx.Then(`^CLI returns exit code 1$`, PurgeDataProductFailed)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, CheckError)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, AssertErrorMessages)
}
