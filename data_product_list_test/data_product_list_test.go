package dataproductlist

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"test-case/testutils"
	"testing"
	"fmt"

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

func CreateDataProductCommand(dataProduct string, description string, enabled string) error {
	dataProduct = ut.ProcessString(dataProduct)
	commandString := "../gravity-cli product create "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	if description != testutils.IgnoreString {
		description := ut.ProcessString(description)
		commandString += " --desc " + description
	}

	commandString += " --schema ./assets/schema.json" 

	if enabled != testutils.IgnoreString {
		if enabled == testutils.TrueString {
			commandString += " --enabled"
		}
	}	
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
	return errors.New("Cli回傳訊息錯誤")
}

func AddRulesetCommand(dataProduct string, RulesetAmount int, ruleset string, method string, event string) error {
	dataProduct = ut.ProcessString(dataProduct)
	for i := 0; i < int(RulesetAmount); i++ {
		cmd := exec.Command(ut.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", "test", "--method", "create", " --schema ./assets/schema.json", "-s", testUtils.Config.JetstreamURL)
		err = cmd.Run()
		if err != nil {
			return errors.New("ruleset add 參數錯誤")
		}
	}
	return nil
}

func AddRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("ruleset 創建應該要成功")
}

func ProductListCommand() error {
	cmd := exec.Command(ut.GravityCliString, "product", "list")
	return cmd.Run()
}

func ProductListCommandSuccess(ProductName string, Description string, Enabled string, RulesetAmount string, EventAmount string) error {
	outStr := ut.CmdResult.Stdout
	outStrList = strings.Split(outStr, "\n")
	product = outStrList[2]
	if product[0] == ProductName && product[1] == Description && product[2] == Enabled && product[3] == RulesetAmount && product[4] == EventAmount{
		return nil
	}
	return errors.New("Cli回傳訊息錯誤")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.When(`^創建 data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'"$`, CreateDataProductCommand)
	ctx.Then(`^Cli 回傳 "'(.*?)'" 建立成功$`, CreateDataProductCommandSuccess)
	ctx.When(`^"'(.*?)'" 創建 "'(.*?)'" 個 ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.When(`^使用gravity-cli 列出所有 data product$`, ProductListCommand)
	ctx.Then(`^Cli 回傳 data product ProductName = "'(.*?)'", Description = "'(.*?)'", Enabled="'(.*?)'", RulesetAmount="'(.*?)'", EventAmount="'(.*?)'"$`, ProductListCommandSuccess)
}