package dataproductrulesetupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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
}

type RuleMap map[string]Rule

type JSONData struct {
	Name  string  `json:"name"`
	Rules RuleMap `json:"rules"`
}

var jsonData JSONData
var ut = testutils.TestUtils{}
var originRuleStr string

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

func SaveRuleset(dataProduct string, ruleset string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	var originJSON JSONData 
	err = json.Unmarshal(entry.Value(), &originJSON)
	if err != nil {
		log.Fatal(err)
	}

	originRuleByte, _ := json.Marshal(originJSON.Rules[ruleset])
	originRuleStr = string(originRuleByte)
	return nil
}

func UpdateRulesetCommand(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string, enabled string) error {
	commandString := "../gravity-cli product ruleset update "
	if dataProduct != "[null]" {
		commandString += " " + dataProduct
	}
	if ruleset != "[null]" {
		commandString += " " + ruleset
	}
	if method != "[ignore]" {
		commandString += " --method " + method
	}
	if event != "[ignore]" {
		commandString += " --event " + event
	}
	if pk != "[ignore]" {
		commandString += " --pk " + pk
	}
	if desc != "[ignore]" {
		commandString += " --desc " + desc
	}
	if enabled != "[ignore]" {
		if enabled == "[true]" {
			commandString += " --enabled"
		} else {
			return errors.New("enabled不允許true或ignore以外的值")
		}
	}
	if handler != "[ignore]" {
		commandString += " --handler " + handler
	}
	if schema != "[ignore]" {
		commandString += " --schema " + schema
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return errors.New("ruleset 更改應該要成功")
}

func UpdateRulesetCommandFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("ruleset 更改應該要失敗")
}

func ValidateRulesetUpdate(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string, enable string) error {
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

	if method != "[ignore]" {
		regexMethod := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(method)[1]
		if jsonData.Rules[ruleset].Method != regexMethod {
			return errors.New("method更改失敗")
		}
	}
	if event != "[ignore]" {
		regexEvent := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(event)[1]
		if jsonData.Rules[ruleset].Event != regexEvent {
			return errors.New("event更改失敗")
		}
	}
	if desc != "[ignore]" {
		regexDesc := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(desc)[1]
		if jsonData.Rules[ruleset].Desc != regexDesc {
			return errors.New("desc更改失敗")
		}
	}
	if pk != "[ignore]" {
		regexPk := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(pk)[1]
		expectedPK := strings.Join(jsonData.Rules[ruleset].PrimaryKey, ",")
		if expectedPK != regexPk {
			return errors.New("pk更改失敗")
		}
	}
	if handler != "[ignore]" {
		regexHandler := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(handler)[1]
		fileContent, err := os.ReadFile(regexHandler)
		if err != nil {
			return errors.New("使用nats驗證時 handler.js 開啟失敗")
		}
		rulesetHandler := jsonData.Rules[ruleset].Handler.(map[string]interface{})
		handlerScript := rulesetHandler["script"].(string)
		if string(fileContent) != handlerScript {
			return errors.New("handler更改失敗")
		}
	}
	if schema != "[ignore]" {
		regexSchema := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(schema)[1]
		fileContent, err := os.ReadFile(regexSchema)
		if err != nil {
			return errors.New("使用nats驗證時 schema.json 開啟失敗")
		}
		natsSchema, _ := json.Marshal(jsonData.Rules[ruleset].Schema)
		var fileJSON interface{}
		err = json.Unmarshal(fileContent, &fileJSON)
		if err != nil {
			return errors.New("使用nats驗證時 schema.json 解碼失敗")
		}
		fileSchemaByte, _ := json.Marshal(fileJSON)
		fileSchema := strings.Join(strings.Fields(string(fileSchemaByte)), "")
		if fileSchema != string(natsSchema) {
			return errors.New("schema更改失敗")
		}
	}
	return nil
}

func ValidateRulesetNotChange(dataProduct string, ruleset string) error {
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	err = json.Unmarshal(entry.Value(), &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	ruleByte, _ := json.Marshal(jsonData.Rules[ruleset])
	if string(ruleByte) != originRuleStr {
		return errors.New("ruleset資料有異動")
	}
	return nil
}

// func AssertErrorMessages(expected string) error {
// 	Todo
// 	return nil
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 data product "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'"$`, ut.CreateDataProductRuleset)
	ctx.Given(`^儲存 nats 現有 data product ruleset 副本 "'(.*?)'" "'(.*?)'"$`, SaveRuleset)
	ctx.When(`^更新 dataproduct "'(.*?)'" ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, UpdateRulesetCommand)
	ctx.Then(`^Cli 回傳更改成功$`, UpdateRulesetCommandSuccess)
	ctx.Then(`^Cli 回傳更改失敗$`, UpdateRulesetCommandFailed)
	ctx.Then(`^使用 nats jetstream查詢 "'(.*?)'" 的 "'(.*?)'" 參數更改成功 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, ValidateRulesetUpdate)
	// ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
	ctx.Then(`^使用 nats jetstream 查詢 data product "'(.*?)'" 的 "'(.*?)'" 資料無任何改動$`, ValidateRulesetNotChange)
}
