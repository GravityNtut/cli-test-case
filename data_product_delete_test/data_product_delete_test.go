package dataproductdelete

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

func SearchDataProductByCLIFail(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return errors.New("Cli 預期不存在，但查詢到了")
}

func SearchDataProductByJetstreamFail(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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
	commandString := "../gravity-cli product delete "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDataProductSuccess(productName string) error {
	outStr := ut.CmdResult.Stdout
	if outStr == "Product \""+productName+"\" was deleted\n" {
		return nil
	}
	return errors.New("Cli回傳訊息錯誤")
}

func DeleteDataProductFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("data product 刪除應該要失敗")
}

// TODO
// func AssertErrorMessages(errorMessage string) error {
// 	outErr := ut.CmdResult.Stderr
// 	if outErr == errorMessage {
// 		return nil
// 	}
// 	return errors.New("Cli回傳訊息錯誤")
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.RestartDocker()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 data product "'(.*?)'" "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.When(`^刪除 data product "'(.*?)'"$`, DeleteDataProductCommand)
	ctx.Then(`^Cli 回傳 "'(.*?)'" 刪除成功$`, DeleteDataProductSuccess)
	ctx.Then(`^Cli 回傳刪除失敗$`, DeleteDataProductFail)
	ctx.Then(`^使用 gravity-cli 查詢 "'(.*?)'" 不存在$`, SearchDataProductByCLIFail)
	ctx.Then(`^使用 nats jetstream 查詢 "'(.*?)'" 不存在$`, SearchDataProductByJetstreamFail)
	// ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
}
