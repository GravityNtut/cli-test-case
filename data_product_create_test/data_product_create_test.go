package data_product_create

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
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

func CreateDataProductCommand(dataProduct string, description string, schema string) error {
	dataProduct = ut.ProcessString(dataProduct)
	commandString := "../gravity-cli product create "
	if dataProduct != "[null]" {
		commandString += dataProduct
	}
	if description != "[ignore]" {
		if description == "[null]" {
			commandString += " --desc"
		} else {
			description := ut.ProcessString(description)
			commandString += " --desc \"" + description + "\""
		}
	}

	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " --enabled" + " -s " + ut.Config.JetstreamURL
	ut.ExecuteShell(commandString)

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

func CreateDataProductCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("Cli回傳訊息錯誤")
}

func SearchDataProductByCLISuccess(dataProduct string) error {
	dataProduct = ut.ProcessString(dataProduct)
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	return err
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	dataProduct = ut.ProcessString(dataProduct)
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return nil
		}
	}
	return errors.New("jetstream裡未創建成功")
}

func AssertErrorMessages(errorMessage string) error {
	// TODO
	// outErr := ut.CmdResult.stderr
	// if outErr == errorMessage {
	// 	return nil
	// }
	// return errors.New("Cli回傳訊息錯誤")
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, ut.CheckDispatcherService)
	ctx.When(`^創建一個data product "([^"]*)" 註解 "([^"]*)" schema檔案 "([^"]*)"$`, CreateDataProductCommand)
	ctx.Step(`^Cli回傳"([^"]*)"建立成功$`, CreateDataProductCommandSuccess)
	ctx.Step(`^Cli回傳建立失敗$`, CreateDataProductCommandFail)
	ctx.Step(`^使用gravity-cli查詢data product 列表 "([^"]*)" 存在$`, SearchDataProductByCLISuccess)
	ctx.Step(`^使用nats jetstream 查詢 data product 列表 "([^"]*)" 存在$`, SearchDataProductByJetstreamSuccess)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
}
