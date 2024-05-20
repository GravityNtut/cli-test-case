package dataproductlist

import (
	"context"
	"errors"
	"os/exec"
	"test-case/testutils"
	"testing"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

const (
	TrueString       = "[true]"
	GravityCliString = "../gravity-cli"
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

func CreateDataProductCommand(productAmount int,dataProduct string, description string, enabled string) error {
	for i:=0; i<productAmount; i++ {
		var dataProductName string
		if i==0 {
			dataProductName = ut.ProcessString(dataProduct)
		}else {
			dataProductName = dataProduct + "_" + strconv.Itoa(i)
		}
		commandString := "../gravity-cli product create "
		commandString += dataProductName
		if description != testutils.IgnoreString {
			description := ut.ProcessString(description)
			commandString += " --desc " + description
		}

		commandString += " --schema './assets/schema.json'" 

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

func CreateDataProductCommandSuccess(productAmount int, productName string) error {
	outStr := ut.CmdResult.Stdout
	productName = ut.ProcessString(productName)
	productAmount-=1
	if( productAmount>1){
		if outStr == "Product \"" + productName + "_" + strconv.Itoa(productAmount) + "\" was created\n" {
			return nil
		}
	}else{
		if outStr == "Product \"" + productName + "\" was created\n" {
			return nil
		}
	}
	
	// return errors.New("Cli回傳訊息create錯誤")
	return errors.New(outStr)
}

func AddRulesetCommand(dataProduct string, RulesetAmount int) error {
	dataProduct = ut.ProcessString(dataProduct)
	for i := 0; i < int(RulesetAmount); i++ {
		ruleset := dataProduct + "Created"
		if(i!=0){
			ruleset += strconv.Itoa(i)
		}
		event := "--event="+ruleset
		cmd := exec.Command(GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--enabled", event, "--method=create", "--schema='./assets/schema.json'", "--pk=id")
		cmdString := cmd.String()
		err := ut.ExecuteShell(cmdString)
		// err := cmd.Run()
		if err != nil {
		// 	return errors.New("ruleset add 參數錯誤")
			outStr := ut.CmdResult.Stderr
			return errors.New(outStr)
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

// TODO: publish是否成功都會return 0, 要怎麼判斷是否失敗
func PublishProductEvent(eventAmount int) error {
	for i := 0; i < eventAmount; i++ {
		event := "drinkCreated"
		payload := fmt.Sprintf(`{"id":%d, "name":"test%d", "kcal":0, "price":0}`, i, i)
		cmd := exec.Command(GravityCliString, "pub", event, payload)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	return nil
}

func PublishProductEventSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("publish event failed")
}


func ProductListCommand() error {
	cmd := "../gravity-cli product list"
	return ut.ExecuteShell(cmd)
}

func ProductListCommandSuccess(ProductAmount int, dataProduct string, Description string, Enabled string, RulesetAmount string, EventAmount string) error {
	outStr := ut.CmdResult.Stdout
	// return errors.New(outStr);

	outStrList := strings.Split(outStr, "\n")
	// fmt.Println("len:",len(outStrList))
	if(len(outStrList) != ProductAmount+3){
		return errors.New("Cli回傳訊息ProductAmount錯誤")
	}

	if strings.Compare(Enabled, TrueString)==0{ // || strings.Compare(Enabled, "enabled")==0
		Enabled = "enabled"	
	} else{
		Enabled = "disabled"
	}

	if(ProductAmount > 1 ){
		for i := 1; i < ProductAmount; i++ {
			product := outStrList[2+i-1]
			dataProductName := dataProduct + "_" + strconv.Itoa(i)
			if !(strings.Contains(product, dataProductName)){
				// return errors.New(outStr)
				return errors.New("Cli回傳list ProductName錯誤")
			}
			if( Description!=blankString1 && Description!=blankString2 ){
				if !( strings.Contains(product, Description)){
					return errors.New("Cli回傳list Description錯誤")
				}
			}
			if !(strings.Contains(product, Enabled)){
				return errors.New("Cli回傳list Enabled錯誤")
			}
		}
	}else{
		product := outStrList[2+ProductAmount-1]
		// outStrList := strings.Split(product, " ")
		outStrList = strings.Fields(product)
		index := 0
		if( Description!=blankString1 && Description!=blankString2 ){
			index++
		}
		if ( outStrList[2+index] != RulesetAmount ){
			// return errors.New("Cli回傳list RulesetAmount錯誤")
			// return errors.New("Cli回傳list RulesetAmount錯誤" + outStrList[0])
			// return errors.New(product)
			return errors.New("Cli回傳list RulesetAmount錯誤" + outStrList[0] + "/" + outStrList[1] + "/" +  outStrList[2] + "/" +  outStrList[3])
		}
		if ( outStrList[3+index] != EventAmount ){
			// return errors.New("Cli回傳list EventAmount錯誤")
			return errors.New("Cli回傳list EventAmount錯誤" + outStrList[3+index])
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
	ctx.Then(`^Cli 回傳第 "'(.*?)'" 個 "'(.*?)'" 建立成功$`, CreateDataProductCommandSuccess)
	ctx.When(`^對"'(.*?)'" 創建 "'(.*?)'" 個 ruleset$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.When(`^對Event做 "'(.*?)'" 次 publish$`, PublishProductEvent)
	ctx.Then(`^publish 成功$`, PublishProductEventSuccess)
	ctx.When(`^使用gravity-cli 列出所有 data product$`, ProductListCommand)
	ctx.Then(`^回傳 data product ProductAmount = "'(.*?)'", ProductName = "'(.*?)'", Description = "'(.*?)'", Enabled="'(.*?)'", RulesetAmount="'(.*?)'", EventAmount="'(.*?)'"$`, ProductListCommandSuccess)
}