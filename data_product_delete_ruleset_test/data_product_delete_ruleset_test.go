package dataproductrulesetdelete

import (
	"context"
	"fmt"
	"os/exec"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
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

func DeleteRulesetCommand(productName string, rulesetName string) error {
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
	return fmt.Errorf("ruleset 刪除應該要失敗")
}

func DeleteRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("ruleset 刪除應該要成功")
}

func SearchRulesetByCLINotExists(dataProduct string, ruleset string) error {
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return fmt.Errorf("ruleset 應該不存在")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 date product "'(.*?)'" enabled "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'" enabled "'(.*?)'"$`, ut.CreateDataProductRuleset)
	ctx.When(`^刪除 "'(.*?)'" 的 ruleset "'(.*?)'"$`, DeleteRulesetCommand)
	ctx.Then(`^Cli 回傳刪除失敗$`, DeleteRulesetCommandFailed)
	ctx.Then(`^Cli 回傳刪除成功$`, DeleteRulesetCommandSuccess)
	ctx.Then(`^使用 gravity-cli 查詢 "'(.*?)'" 的 "'(.*?)'" 不存在$`, SearchRulesetByCLINotExists)
}
