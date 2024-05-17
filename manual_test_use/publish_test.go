package dataproductrulesetadd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"test-case/testutils"
	"testing"

	gravity_sdk_types_product_event "github.com/BrobridgeOrg/gravity-sdk/v2/types/product_event"
	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
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

type JSONData struct {
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

var jsonData JSONData
var EventCount = 3

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

func AddRulesetCommand(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string) error {
	dataProduct = ut.ProcessString(dataProduct)
	ruleset = ut.ProcessString(ruleset)
	commandString := "../gravity-cli product ruleset add "
	if dataProduct != testutils.NullString {
		commandString += " " + dataProduct
	}
	if ruleset != testutils.NullString {
		commandString += " " + ruleset
	}
	if event != testutils.IgnoreString {
		event := ut.ProcessString(event)
		commandString += " --event " + event
	}
	if method != testutils.IgnoreString {
		method := ut.ProcessString(method)
		commandString += " --method " + method
	}
	if pk != testutils.IgnoreString {
		pk := ut.ProcessString(pk)
		commandString += " --pk " + pk
	}
	if desc != testutils.IgnoreString {
		desc := ut.ProcessString(desc)
		commandString += " --desc " + desc
	}
	if handler != testutils.IgnoreString {
		commandString += " --handler " + handler
	}
	if schema != testutils.IgnoreString {
		commandString += " --schema " + schema
	}
	commandString += " --enabled -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
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
	ruleset = ut.ProcessString(ruleset)
	method = ut.ProcessString(method)
	event = ut.ProcessString(event)
	desc = ut.ProcessString(desc)
	pk = ut.ProcessString(pk)

	if ruleset != testutils.NullString {
		if ruleset != jsonData.Rules[ruleset].Name {
			return errors.New("ruleset 與 nats 資訊不符")
		}
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
	return nil
}

var DataProduct string
var Ruleset string
var Event string

const Method string = "create"
const Pk string = "id,name"
const RulesetDesc string = "drink創建事件"
const Handler string = "./assets/handler.js"
const Schema string = "./assets/schema.json"
const Enabled string = "true"

func CreateDataProduct(dataProduct string) error {
	decription := "drink資料"
	schema := "./assets/schema.json"
	enabled := "true"
	cmd := exec.Command(testutils.GravityCliString, "product", "create", dataProduct, "--desc", decription, "--schema", schema, "--enabled="+enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("create data product failed")
	}
	return nil
}

func CreateDataProductRuleset(dataProduct string, ruleset string) error {
	DataProduct = dataProduct
	Ruleset = ruleset
	Event = ruleset
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", DataProduct, Ruleset, "--event", Event, "--method", Method, "--desc", RulesetDesc, "--pk", Pk, "--handler", Handler, "--schema", Schema, "--enabled="+Enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
	}
	return nil
}

func PublishProductEvent(numbersOfEvents int) error {
	EventCount = numbersOfEvents
	for i := 0; i < EventCount; i++ {
		payload := `{"id":%d, "name":"test%d", "kcal":%d, "price":%d}`
		result := fmt.Sprintf(payload, i+1, i+1, i*100, i+20)
		cmd := exec.Command(testutils.GravityCliString, "pub", Event, result, "-s", ut.Config.JetstreamURL)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func DisplayData(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	var pe gravity_sdk_types_product_event.ProductEvent

	ch := make(chan *nats.Msg, EventCount)
	js.ChanSubscribe("$GVT.default.DP."+dataProduct+".*.EVENT.>", ch)
	for i := 0; i < EventCount; i++ {
		msg := <-ch 
		err = proto.Unmarshal(msg.Data, &pe)
		if err != nil {
			fmt.Printf("Failed to parsing product event: %v", err)
		}

		r, err := pe.GetContent()
		if err != nil {
			fmt.Printf("Failed to parsing content: %v", err)
		}
		fmt.Println(pe.EventName)
		fmt.Println(r.AsMap())
	}
	return nil
}

// var productSubscriberName string = ""
// var productSubscriberPartitions []int = []int{-1}
// var productSubscriberStartSeq uint64 = 1

// type ProductCommandContext struct {
// 	Config    *configs.Config
// 	Logger    *zap.Logger
// 	Connector *connector.Connector
// 	Product   *product.Product
// 	Cmd       *cobra.Command
// 	Args      []string
// }

// var cctx = &ProductCommandContext{}

// func DisplayData(productName string) error {

// 	opts := subscriber_sdk.NewOptions()
// 	opts.Verbose = true
// 	s := subscriber_sdk.NewSubscriberWithClient(productSubscriberName, cctx.Connector.GetClient(), opts)
// 	_, err := s.Subscribe(productName, func(msg *nats.Msg) {

// 		var pe gravity_sdk_types_product_event.ProductEvent

// 		err := proto.Unmarshal(msg.Data, &pe)
// 		if err != nil {
// 			fmt.Printf("Failed to parsing product event: %v", err)
// 			msg.Ack()
// 			return
// 		}

// 		md, _ := msg.Metadata()

// 		r, err := pe.GetContent()
// 		if err != nil {
// 			fmt.Printf("Failed to parsing content: %v", err)
// 			msg.Ack()
// 			return
// 		}

// 		// Convert data to JSON
// 		event := map[string]interface{}{
// 			"header":     msg.Header,
// 			"subject":    msg.Subject,
// 			"seq":        md.Sequence.Consumer,
// 			"timestamp":  md.Timestamp,
// 			"product":    productName,
// 			"event":      pe.EventName,
// 			"method":     pe.Method.String(),
// 			"table":      pe.Table,
// 			"primaryKey": pe.PrimaryKeys,
// 			"payload":    r.AsMap(),
// 		}

// 		data, _ := json.MarshalIndent(event, "", "  ")
// 		fmt.Println(string(data))
// 		msg.Ack()
// 	}, subscriber_sdk.Partition(productSubscriberPartitions...), subscriber_sdk.StartSequence(productSubscriberStartSeq))
// 	if err != nil {
// 		cctx.Cmd.SilenceUsage = true
// 		return err
// 	}

// 	select {}
// }

// func AssertErrorMessages(expected string) error {
// 	if cmdResult.Stderr == expected {
// 		return nil
// 	}
// TODO
// 	return fmt.Errorf("應有錯誤訊息: %s", expected)
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^已有data product "'(.*?)'"$`, CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	ctx.When(`^"'(.*?)'" 創建ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建失敗$`, AddRulesetCommandFailed)
	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.Then(`^已 publish "'(.*?)'" 筆 Event$`, PublishProductEvent)
	ctx.Then(`^顯示 "'(.*?)'" 資料$`, DisplayData)
	ctx.Then(`^使用gravity-cli 查詢 "'(.*?)'" 的 "'(.*?)'" 存在$`, SearchRulesetByCLISuccess)
	ctx.Then(`使用nats jetstream 查詢 "'(.*?)'" 的 "'(.*?)'" 存在，且參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" 正確$`, SearchRulesetByNatsSuccess)
	// ctx.Then(`^應有錯誤訊息 "'(.*?)'"$`, AssertErrorMessages)
}
