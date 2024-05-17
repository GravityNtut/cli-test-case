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

func CreateDataProductCommand(productAmount int,dataProduct string, description string, enabled string) error {
	for i=0; i<productAmount; i++ {
		if i==0 {
			dataProduct = ut.ProcessString(dataProduct)
		}else {
			dataProduct = dataProduct + "_" + strconv.Itoa(i)
		}
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

func PublishProductEvent(event string, eventAmount int) error {
	for i := 0; i < EventCount; i++ {
		payload := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":0, "price":0}`, i, i)
		cmd := exec.Command(testutils.GravityCliString, "pub", event, payload)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	payload = ""
	return nil
}

func PublishProductEventSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("publish event failed")
}


func ProductListCommand() error {
	cmd := exec.Command(ut.GravityCliString, "product", "list")
	return cmd.Run()
}

func ProductListCommandSuccess(ProductAmount int, ProductName string, Description string, Enabled string, RulesetAmount string, EventAmount string) error {
	outStr := ut.CmdResult.Stdout
	outStrList = strings.Split(outStr, "\n")
	for i := 0; i < ProductAmount; i++ {
		product = outStrList[2 + i]
		if i==0 {
			dataProduct = ut.ProcessString(ProductName)
		}else {
			dataProduct = ProductName + "_" + strconv.Itoa(i)
		}


		if enabled == TrueString {
			enabled = "enabled"	
		} else{
			enabled = "disabled"
		}

		if !(strings.Contains(product, dataProduct) && strings.Contains(product, Description) && strings.Contains(product, Enabled) && strings.Contains(product, RulesetAmount) && strings.Contains(product, EventAmount)) {
			return errors.New("Cli回傳訊息錯誤")
		}
	}
	
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.When(`^創建 "'(.*?)'" 個 data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'"$`, CreateDataProductCommand)
	ctx.Then(`^Cli 回傳 "'(.*?)'" 建立成功$`, CreateDataProductCommandSuccess)
	ctx.When(`^"'(.*?)'" 創建 "'(.*?)'" 個 ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.When(`^對 "'(.*?)'" 做 "'(.*?)'" 次 publish$`, PublishProductEvent)
	ctx.Then(`^publish 成功$`, PublishProductEventSuccess)
	ctx.When(`^使用gravity-cli 列出所有 data product$`, ProductListCommand)
	ctx.Then(`^Cli 回傳 data product ProductAmount = "'(.*?)'", ProductName = "'(.*?)'", Description = "'(.*?)'", Enabled="'(.*?)'", RulesetAmount="'(.*?)'", EventAmount="'(.*?)'"$`, ProductListCommandSuccess)
}