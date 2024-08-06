package dataproductrulesetdelete

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
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

func DeleteRulesetCommand(rulesetName string, productName string) error {
	commandString := "../gravity-cli product ruleset delete "
	if productName != testutils.NullString {
		commandString += " " + productName
	}
	if rulesetName != testutils.NullString {
		commandString += " " + rulesetName
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRulesetCommandFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return fmt.Errorf("the data product ruleset deletion should fail")
}

func DeleteRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("the data product ruleset deletion should succeed")
}

func SearchRulesetByCLINotExists(ruleset string, dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return fmt.Errorf("ruleset should not exist")
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
	ctx.When(`^Delete ruleset "'(.*?)'" for data product "'(.*?)'"$`, DeleteRulesetCommand)
	ctx.Then(`^CLI returns exit code 1$`, DeleteRulesetCommandFailed)
	ctx.Then(`^CLI returned successfully deleted$`, DeleteRulesetCommandSuccess)
	ctx.Then(`^Using gravity-cli to query that "'(.*?)'" does not exist for "'(.*?)'"$`, SearchRulesetByCLINotExists)
}
