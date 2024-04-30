package dataproductrulesetupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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
	if dataProduct != testutils.NullString {
		commandString += " " + dataProduct
	}
	if ruleset != testutils.NullString {
		commandString += " " + ruleset
	}
	if method != testutils.IgnoreString {
		commandString += " --method " + method
	}
	if event != testutils.IgnoreString {
		commandString += " --event " + event
	}
	if pk != testutils.IgnoreString {
		commandString += " --pk " + pk
	}
	if desc != testutils.IgnoreString {
		commandString += " --desc " + desc
	}
	if enabled == testutils.TrueString {
		commandString += " --enabled"
	} else if enabled == testutils.FalseString {
		commandString += " --enabled=false"
	} else if enabled != testutils.IgnoreString {
		return errors.New("enabled 參數錯誤")
	}
	if handler != testutils.IgnoreString {
		commandString += " --handler " + handler
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

func ValidateRulesetUpdate(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string, enabled string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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

	if err := ut.ValidateField(jsonData.Rules[ruleset].Method, method); err != nil {
		return err
	}
	if err := ut.ValidateField(jsonData.Rules[ruleset].Event, event); err != nil {
		return err
	}
	if err := ut.ValidateField(jsonData.Rules[ruleset].Desc, desc); err != nil {
		return err
	}
	pkStr := strings.Join(jsonData.Rules[ruleset].PrimaryKey, ",")
	if err := ut.ValidateField(pkStr, pk); err != nil {
		return err
	}

	if err := ut.ValidateHandler(jsonData.Rules[ruleset].Handler, handler); err != nil {
		return err
	}

	if err := ut.ValidateSchema(jsonData.Rules[ruleset].Schema, schema); err != nil {
		return err
	}

	if err := ut.ValidateEnabled(jsonData.Rules[ruleset].Enabled, enabled); err != nil {
		return err
	}
	return nil
}

func ValidateRulesetNotChange(dataProduct string, ruleset string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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

// TODO
// func AssertErrorMessages(expected string) error {
// return nil
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有 data product "'(.*?)'" enabled "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'" enabled "'(.*?)'"$`, ut.CreateDataProductRuleset)
	ctx.Given(`^儲存 nats 現有 data product ruleset 副本 "'(.*?)'" "'(.*?)'"$`, SaveRuleset)
	ctx.When(`^更新 dataproduct "'(.*?)'" ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, UpdateRulesetCommand)
	ctx.Then(`^Cli 回傳更改成功$`, UpdateRulesetCommandSuccess)
	ctx.Then(`^Cli 回傳更改失敗$`, UpdateRulesetCommandFailed)
	ctx.Then(`^使用 nats jetstream查詢 "'(.*?)'" 的 "'(.*?)'" 參數更改成功 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, ValidateRulesetUpdate)
	// ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
	ctx.Then(`^使用 nats jetstream 查詢 data product "'(.*?)'" 的 "'(.*?)'" 資料無任何改動$`, ValidateRulesetNotChange)
}
