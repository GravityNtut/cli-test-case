package dataproductlist

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
)

const (
	blankString1 = "\"\""
	blankString2 = "\" \""
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

func CreateDataProductCommand(productAmount int, dataProduct string, description string, enabled string) error {
	var dataProductName string
	dataProductName = ut.ProcessString(dataProduct)
	for i := 0; i < productAmount; i++ {
		if i != 0 {
			dataProductName = dataProduct + "_" + strconv.Itoa(i)
		}
		if description != testutils.IgnoreString {
			description = ut.ProcessString(description)
		}

		enabledString := ""
		if enabled != testutils.IgnoreString {
			if enabled == testutils.TrueString {
				enabledString += " --enabled"
			}
		}
		cmd := exec.Command("../gravity-cli", "product", "create", dataProductName, "--desc", description, "--schema", "./assets/schema.json", enabledString)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddRulesetCommand(dataProduct string, RulesetAmount int) error {
	dataProduct = ut.ProcessString(dataProduct)
	for i := 0; i < RulesetAmount; i++ {
		ruleset := dataProduct + "Created"
		if i != 0 {
			ruleset += strconv.Itoa(i)
		}
		event := "--event " + ruleset
		cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", ruleset, "--enabled", event, "--method", "create", "--schema", "./assets/schema.json", "--pk", "id")
		err := cmd.Run()
		if err != nil {
			return errors.New(cmd.String())
		}
	}
	return nil
}

func PublishProductEvent(eventAmount int) error {
	for i := 0; i < eventAmount; i++ {
		event := "drinkCreated"
		payload := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":0, "price":0}`, i, i)
		cmd := exec.Command(testutils.GravityCliString, "pub", event, payload)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	return nil
}

func ProductListCommand() error {
	const cmd = "../gravity-cli product list"
	return ut.ExecuteShell(cmd)
}

func ProductListCommandSuccess(productAmount int, dataProduct string, description string, enabled string, rulesetAmount string, eventAmount string) error {
	outStr := ut.CmdResult.Stdout

	outStrList := strings.Split(outStr, "\n")
	if len(outStrList) != productAmount+3 {
		return errors.New("Cli回傳訊息ProductAmount錯誤")
	}

	if strings.Compare(enabled, testutils.TrueString) == 0 {
		enabled = "enabled"
	} else {
		enabled = "disabled"
	}

	if productAmount > 1 {
		for i := 1; i < productAmount; i++ {
			product := outStrList[2+i-1]
			dataProductName := dataProduct + "_" + strconv.Itoa(i)
			if !(strings.Contains(product, dataProductName)) {
				return errors.New("Cli回傳list ProductName錯誤")
			}
			if description != blankString1 && description != blankString2 {
				if !(strings.Contains(product, description)) {
					return errors.New("Cli回傳list Description錯誤")
				}
			}
			if !(strings.Contains(product, enabled)) {
				return errors.New("Cli回傳list Enabled錯誤")
			}
		}
	} else {
		product := outStrList[2+productAmount-1]
		productItem := strings.Fields(product)
		index := 0
		if description != blankString1 && description != blankString2 {
			index++
		}
		if productItem[2+index] != rulesetAmount {
			return errors.New("Cli回傳list RulesetAmount錯誤")
		}
		if productItem[3+index] != eventAmount {
			return errors.New("Cli回傳list EventAmount錯誤")
		}
	}

	return nil
}

func ProductListCommandFail() error {
	outStr := ut.CmdResult.Stderr
	if strings.Contains(outStr, "Error: No available products") {
		return nil
	}
	return errors.New(outStr)
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.When(`^創建 "'(.*?)'" 個 data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'"$`, CreateDataProductCommand)
	ctx.When(`^對"'(.*?)'" 創建 "'(.*?)'" 個 ruleset$`, AddRulesetCommand)
	ctx.When(`^對Event做 "'(.*?)'" 次 publish$`, PublishProductEvent)
	ctx.When(`^使用gravity-cli 列出所有 data product$`, ProductListCommand)
	ctx.Then(`^Cli 回傳 "'(.*?)'" 個 data product, 每個 data product 裡面的名字為 "'(.*?)'", 描述內容為 "'(.*?)'", Enabled 的狀態為 "'(.*?)'", Ruleset 的數量為 "'(.*?)'" 個, 以及 Event 總共發布 "'(.*?)'" 個 $`, ProductListCommandSuccess)
	ctx.Then(`^回傳 Error: No available products$`, ProductListCommandFail)
}
