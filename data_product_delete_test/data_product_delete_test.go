package data_product_delete

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

func SearchDataProductByCLIFail(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return errors.New("Cli 預期不存在，但查詢到了")
}

func SearchDataProductByJetstreamFail(dataProduct string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return errors.New("Jetstream 預期不存在，但查詢到了")
		}
	}
	return nil
}

func DeleteDataProductCommand(dataProduct string) error {
	commandString := "../gravity-cli product delete " + dataProduct + " -s " + ut.Config.JetstreamURL
	ut.ExecuteShell(commandString)
	return nil
}

func DeleteDataProductSuccess() error {
	return ut.CmdResult.Err
}

func DeleteDataProductFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("data product 刪除應該要失敗")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有data product "([^"]*)"$`, ut.CreateDataProduct)
	ctx.When(`^刪除data product "([^"]*)"$`, DeleteDataProductCommand)
	ctx.Then(`^data product 刪除成功$`, DeleteDataProductSuccess)
	ctx.Then(`^data product 刪除失敗$`, DeleteDataProductFail)
	ctx.Then(`^使用gravity-cli查詢data product 列表 "([^"]*)" 不存在$`, SearchDataProductByCLIFail)
	ctx.Then(`^使用nats jetstream 查詢 data product 列表 "([^"]*)" 不存在$`, SearchDataProductByJetstreamFail)
}
