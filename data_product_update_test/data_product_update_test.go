package dataproductupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type JSONData struct {
	Name    string      `json:"name"`
	Desc    string      `json:"desc"`
	Enabled bool        `json:"enabled"`
	Schema  interface{} `json:"schema"`
}

var newJSONData JSONData
var ut testutils.TestUtils
var originJSONData string

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
			StopOnFailure: false,
			TestingT:      t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func UpdateDataProductCommand(dataProduct string, description string, enabled string, schema string) error {
	commandString := "../gravity-cli product update "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	if description != testutils.IgnoreString {
		commandString += " --desc " + description
	}

	if enabled == testutils.TrueString {
		commandString += " --enabled"
	} else if enabled == testutils.FalseString {
		commandString += " --enabled=false"
	} else if enabled != testutils.IgnoreString {
		return errors.New("enabled 參數錯誤")
	}

	if schema != testutils.IgnoreString {
		commandString += " --schema " + schema
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDataProductCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("更新資料錯誤")
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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

// func AssertErrorMessages(errorMessage string) error {
// 	outErr := cmdResult.stderr
// 	if outErr == errorMessage {
// 		return nil
// 	}
// 	return errors.New("Cli回傳訊息錯誤")
// }

func DataProductNotChanges(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	if string(entry.Value()) == originJSONData {
		return nil
	}
	return errors.New("與原始資料不符")
}

func DataProductUpdateSuccess(dataProduct string, desc string, schema string, enabled string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &newJSONData)
	if err != nil {
		fmt.Println("解碼 JSON 時出現錯誤:", err)
		return err
	}

	if err := ut.ValidateSchema(newJSONData.Schema, schema); err != nil {
		return err
	}

	if err := ut.ValidateEnabled(newJSONData.Enabled, enabled); err != nil {
		return err
	}

	if err := ut.ValidateField(newJSONData.Desc, desc); err != nil {
		return err
	}

	if dataProduct != newJSONData.Name {
		return errors.New("資料更新失敗")
	}
	return nil
}

func UpdateDataProductCommandSuccess() error {
	if ut.CmdResult.Err != nil {
		return errors.New("更新資料錯誤")
	}
	return nil
}

func StoreNowDataProduct(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	originJSONData = string(entry.Value())
	return nil
}

// TODO
// func AssertErrorMessages(errorMessage string) error {

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 data product "'(.*?)'" "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^儲存 nats 現有 data product 副本 "'(.*?)'"$`, StoreNowDataProduct)
	ctx.When(`^更新 data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, UpdateDataProductCommand)
	ctx.Then(`^Cli 回傳更改成功$`, UpdateDataProductCommandSuccess)
	ctx.Then(`^Cli 回傳更改失敗$`, UpdateDataProductCommandFail)
	// ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
	ctx.Then(`^使用 nats jetstream 查詢 "'(.*?)'" 參數更改成功 "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, DataProductUpdateSuccess)
	ctx.Then(`^使用 nats jetstream 查詢 "'(.*?)'" 參數無任何改動$`, DataProductNotChanges)
}
