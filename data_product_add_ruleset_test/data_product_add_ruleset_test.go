package data_product_ruleset_add

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
	ut.LoadConfig()
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

func AddRulesetCommand(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string) error {
	dataProduct = ut.ProcessString(dataProduct)
	ruleset = ut.ProcessString(ruleset)
	commandString := "../gravity-cli product ruleset add " + dataProduct + " " + ruleset
	if event != "[ignore]" {
		event := ut.ProcessString(event)
		commandString += " --event " + event
	}
	if method != "[ignore]" {
		method := ut.ProcessString(method)
		commandString += " --method " + method
	}
	if pk != "[ignore]" {
		pk := ut.ProcessString(pk)
		commandString += " --pk " + pk
	}
	if desc != "[ignore]" {
		if desc == "[null]" {
			commandString += " --desc "
		} else {
			desc := ut.ProcessString(desc)
			commandString += " --desc \"" + desc + "\""
		}
	}
	if handler != "[ignore]" {
		commandString += " --handler ./assets/" + handler
	}
	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " -s " + ut.Config.JetstreamURL
	ut.ExecuteShell(commandString)
	return nil
}

func AddRulesetCommandFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return fmt.Errorf("ruleset 創建應該要失敗")
}

func AddRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("ruleset 創建應該要成功")
}

func SearchRulesetByCLISuccess(dataProduct string, ruleset string) error {
	dataProduct = ut.ProcessString(dataProduct)
	ruleset = ut.ProcessString(ruleset)
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	return err
}

func AssertErrorMessages(expected string) error {
	// Todo
	// if cmdResult.Stderr == expected {
	// 	return nil
	// }
	// return fmt.Errorf("應有錯誤訊息: %s", expected)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, ut.CheckDispatcherService)

	ctx.Given(`^已有data product "([^"]*)"$`, ut.CreateDataProduct)

	ctx.When(`^"([^"]*)" 創建ruleset "([^"]*)" method "([^"]*)" event "([^"]*)" pk "([^"]*)" desc "([^"]*)" handler "([^"]*)" schema "([^"]*)"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建失敗$`, AddRulesetCommandFailed)

	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.Then(`^使用gravity-cli 查詢 "([^"]*)" 的 "([^"]*)" 成功$`, SearchRulesetByCLISuccess)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
}
