package data_product_update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type JsonData struct {
	Name    string      `json:"name"`
	Desc    string      `json:"desc"`
	Enabled bool        `json:"enabled"`
	Schema  interface{} `json:"schema"`
}

var newJsonData JsonData
var ut testutils.TestUtils
var originJsonData string

func TestFeatures(t *testing.T) {
	ut.LoadConfig()
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

func UpdateDataProductCommand(dataProduct string, description string, enable string, schema string) error {
	commandString := "../gravity-cli product update "
	commandString += dataProduct
	if description != "[ignore]" {
		if description == "[null]" {
			commandString += " --desc"
		} else {
			commandString += " --desc \"" + description + "\""
		}
	}

	if enable != "[ignore]" {
		if enable == "[true]" {
			commandString += " --enabled"
		} else {
			return errors.New("不允許true或ignore以外的輸入")
		}
	}

	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " -s " + ut.Config.JetstreamURL
	ut.ExecuteShell(commandString)

	return nil
}

func UpdateDataProductCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("更新資料錯誤")
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
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
	// outErr := cmdResult.stderr
	// if outErr == errorMessage {
	// 	return nil
	// }
	// return errors.New("Cli回傳訊息錯誤")
	return nil
}

func DataProductNotChanges(dataProduct string, description string, schema string, enabled string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	if string(entry.Value()) == originJsonData {
		return nil
	}
	return errors.New("與原始資料不符")
}

func DataProductUpdateSuccess(dataProduct string, description string, schema string, enabled string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &newJsonData)
	if err != nil {
		fmt.Println("解碼 JSON 時出現錯誤:", err)
		return err
	}

	if schema != "[ignore]" {
		fileContent, err := os.ReadFile("./assets/" + schema)
		if err != nil {
			return err
		}
		schemaString, _ := json.Marshal(newJsonData.Schema)
		var jsonInterface interface{}
		json.Unmarshal(fileContent, &jsonInterface)
		fileSchemaByte, _ := json.Marshal(jsonInterface)
		fileSchema := strings.Join(strings.Fields(string(fileSchemaByte)), "")
		if fileSchema != string(schemaString) {
			return errors.New("schema內容不同")
		}
	}

	var enabledBool bool
	if enabled == "[true]" {
		enabledBool = true
	} else if enabled == "[ignore]" {
		enabledBool = false
	} else {
		return errors.New("不允許true或ignore以外的輸入")
	}

	if description != "[ignore]" {
		if description != newJsonData.Desc {
			return errors.New("description內容不同")
		}
	}

	if dataProduct == newJsonData.Name && enabledBool == newJsonData.Enabled {
		return nil
	}
	return errors.New("資料更新失敗")
}

func UpdateDataProductCommandSuccess() error {
	if ut.CmdResult.Err != nil {
		return errors.New("更新資料錯誤")
	}
	return nil
}

func StoreNowDataProduct(dataProduct string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	originJsonData = string(entry.Value())
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
	ctx.Given(`^儲存nats現有data product "([^"]*)" 副本$`, StoreNowDataProduct)
	ctx.When(`^更新data product "([^"]*)" 註解 "([^"]*)" "([^"]*)" schema檔案 "([^"]*)"$`, UpdateDataProductCommand)
	ctx.Then(`^data product更改成功$`, UpdateDataProductCommandSuccess)
	ctx.Then(`^data product更改失敗$`, UpdateDataProductCommandFail)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
	ctx.Then(`^使用nats驗證data product "([^"]*)" description "([^"]*)" schema檔案 "([^"]*)" enabled "([^"]*)" 更改成功`, DataProductUpdateSuccess)
	ctx.Then(`^使用nats驗證data product "([^"]*)" description "([^"]*)" schema檔案 "([^"]*)" enabled "([^"]*)" 資料無任何改動$`, DataProductNotChanges)
}