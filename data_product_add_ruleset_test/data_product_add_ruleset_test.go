package data_product_ruleset_add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type Rule struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Desc       string      `json:"desc"`
	Event      string      `json:"event"`
	Product    string      `json:"product"`
	Method     string      `json:"method"`
	PrimaryKey []string    `json:"primaryKey"`
	Enabled    bool        `json:"enabled"`
	Handler    interface{} `json:"handler"`
	Schema     interface{} `json:"schema"`
	CreatedAt  string      `json:"createdAt"`
	UpdatedAt  string      `json:"updatedAt"`
}

type RuleMap map[string]Rule

type JsonData struct {
	Name            string      `json:"name"`
	Desc            string      `json:"desc"`
	Enabled         bool        `json:"enabled"`
	Rules           RuleMap     `json:"rules"`
	Schema          interface{} `json:"schema"`
	EnabledSnapshot bool        `json:"enabledSnapshot"`
	Snapshot        interface{} `json:"snapshot"`
	Stream          string      `json:"stream"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
}

var jsonData JsonData

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
	commandString := "../gravity-cli product ruleset add "
	if dataProduct != "[null]" {
		commandString += " " + dataProduct
	}
	if ruleset != "[null]" {
		commandString += " " + ruleset
	}
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
		desc := ut.ProcessString(desc)
		commandString += " --desc " + desc
	}
	if handler != "[ignore]" {
		commandString += " --handler " + handler
	}
	if schema != "[ignore]" {
		commandString += " --schema " + schema
	}
	commandString += " --enabled -s " + ut.Config.JetstreamURL
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

func SearchRulesetByNatsSuccess(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &jsonData)
	if err != nil {
		fmt.Println("解碼 JSON 時出現錯誤:", err)
		return err
	}
	ruleset = ut.ProcessString(ruleset)
	//以下四個參數shell可以加上雙引號，因此這裡要進行移除
	method = regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(method)[1]
	event = regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(event)[1]
	desc = regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(desc)[1]
	pk = regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(pk)[1]
	method = ut.ProcessString(method)
	event = ut.ProcessString(event)
	desc = ut.ProcessString(desc)
	pk = ut.ProcessString(pk)

	if ruleset != "[null]" {
		if ruleset != jsonData.Rules[ruleset].Name {
			return errors.New("NATS 查詢 ruleset 名稱不正確")
		}
	}

	if method != "[ignore]" {
		if method != jsonData.Rules[ruleset].Method {
			return errors.New("NATS 查詢 ruleset method資訊不正確")
		}
	}

	if event != "[ignore]" {
		if event != jsonData.Rules[ruleset].Event {
			return errors.New("NATS 查詢 ruleset event資訊不正確")
		}
	}

	if desc != "[ignore]" {
		if desc != jsonData.Rules[ruleset].Desc {
			return errors.New("NATS 查詢 ruleset desc資訊不正確")
		}
	}


	if pk != "[ignore]" {
		expectedPK := strings.Join(jsonData.Rules[ruleset].PrimaryKey, ",")
		if pk != expectedPK {
			return errors.New("NATS 查詢 ruleset PK資訊不正確")
		}
	}

	if handler != "[ignore]" {
		regexHandler := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(handler)[1]
		fileContent, err := os.ReadFile(regexHandler)
		if err != nil {
			return errors.New("NATS 查詢 handler.js 開啟失敗")
		}
		rulesetHandler := jsonData.Rules[ruleset].Handler.(map[string]interface{})
		handlerScript := rulesetHandler["script"].(string)
		if string(fileContent) != handlerScript {
			return errors.New("NATS 查詢 ruleset handler.js 資訊不正確")
		}
	}

	if schema != "[ignore]" {
		regexSchema := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(schema)[1]
		fileContent, err := os.ReadFile(regexSchema)
		if err != nil {
			return errors.New("NATS 查詢 schema.json 開啟失敗")
		}
		natsSchema, _ := json.Marshal(jsonData.Rules[ruleset].Schema)
		var fileJson interface{}
		json.Unmarshal(fileContent, &fileJson)
		fileSchemaByte, _ := json.Marshal(fileJson)
		fileSchema := strings.Join(strings.Fields(string(fileSchemaByte)), "")
		if fileSchema != string(natsSchema) {
			return errors.New("NATS 查詢 ruleset schema.json 資訊不正確")
		}
	}
	return nil
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

	ctx.Given(`^已有data product "'(.*?)'"$`, ut.CreateDataProduct)

	ctx.When(`^"'(.*?)'" 創建ruleset "'(.*?)'" method "'(.*?)'" event "'(.*?)'" pk "'(.*?)'" desc "'(.*?)'" handler "'(.*?)'" schema "'(.*?)'"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建失敗$`, AddRulesetCommandFailed)
	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.Then(`^使用gravity-cli 查詢 "'(.*?)'" 的 "'(.*?)'" 存在$`, SearchRulesetByCLISuccess)
	ctx.Then(`使用nats jetstream 查詢 "'(.*?)'" 的 "'(.*?)'" 存在，且參數 method "'(.*?)'" event "'(.*?)'" pk "'(.*?)'" desc "'(.*?)'" handler "'(.*?)'" schema "'(.*?)'" 正確$`, SearchRulesetByNatsSuccess)
	ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
}
