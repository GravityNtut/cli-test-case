package dataproductpublish

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"test-case/testutils"
	"testing"
	"time"

	gravity_sdk_types_product_event "github.com/BrobridgeOrg/gravity-sdk/v2/types/product_event"
	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
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

func CreateDataProductRuleset(dataProduct string, ruleset string, RSMethod string, event string, RSPk string, RSHandler string, RSSchema string, RSEnabled string) error {

	if RSEnabled == testutils.TrueString {
		RSEnabled = "true"
	} else if RSEnabled == testutils.FalseString {
		RSEnabled = "false"
	} else {
		return errors.New("Enable 必須要[true] 或 [false]")
	}
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", event, "--method", RSMethod, "--pk", RSPk, "--handler", RSHandler, "--schema", RSSchema, "--enabled="+RSEnabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
	}
	return nil
}

func GeneratePayloadWithLargeNum(number int) string {
	var payload string
	for i := 0; i < number; i++ {
		payload += "a"
	}
	result := fmt.Sprintf(`{"id":1, "name":"%s", "kcal":1000, "price":100}`, payload)
	return result
}

func PublishEventCommand(event string, payload string) error {
	payload = ut.ProcessString(payload)
	pubString := "../gravity-cli pub " + event + " " + payload
	if err := ut.ExecuteShell(pubString); err != nil {
		return err
	}
	return nil
}

type JSONData struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

func QueryJetstreamEventExist(event string, payload string) error {
	// 移除最外邊的單引號
	payload = regexp.MustCompile(`'?([^']*)'?`).FindStringSubmatch(payload)[1]
	payload = ut.ProcessString(payload)
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *nats.Msg, 1)
	if _, err := js.ChanSubscribe("$GVT.default.EVENT.*", ch); err != nil {
		return fmt.Errorf("js get subscribe failed: %v", err)
	}

	m := <-ch
	var jsonData JSONData
	if err := json.Unmarshal(m.Data, &jsonData); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}
	fmt.Println("Event: ", jsonData.Event)
	result, _ := Base64ToString(jsonData.Payload)
	fmt.Println("Payload: ", result)

	if jsonData.Event != event {
		return fmt.Errorf("預期的event: %s, 實際的event: %s", event, jsonData.Event)
	}
	if result != payload {
		return fmt.Errorf("預期的payload: %s, 實際的payload: %s", payload, result)
	}
	return nil
}

func CheckDPStreamDPNotExist(dataProduct string) error {
	const EventCount = 1
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *nats.Msg, EventCount)
	sub, err := js.ChanSubscribe("$GVT.default.DP."+dataProduct+".*.EVENT.>", ch)
	if err != nil {
		return fmt.Errorf("subscribe failed %s", err.Error())
	}
	time.Sleep(1 * time.Second)
	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("unsubscribe failed %s", err.Error())
	}
	if len(ch) != 0 {
		return fmt.Errorf("預期不會進到GVT_default_DP裡，但是進了")
	}
	return nil
}

func Base64ToString(base64Str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func UpdateDataProductCommand(dataProduct string) error {
	cmd := exec.Command(testutils.GravityCliString, "product", "update", dataProduct, "--enabled=true", "-s", ut.Config.JetstreamURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	fmt.Println(stderr.String())
	if err != nil {
		return errors.New("data product update failed")
	}
	return nil
}

func UpdateRulesetCommand(dataProduct string, ruleset string) error {
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "update", dataProduct, ruleset, "--enabled=true", "-s", ut.Config.JetstreamURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	fmt.Println(stderr.String())
	if err != nil {
		return errors.New("ruleset update failed")
	}
	return nil
}

func CheckDPStreamDPExist(dataProduct string, event string, payload string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	var pe gravity_sdk_types_product_event.ProductEvent

	ch := make(chan *nats.Msg, 1)
	sub, err := js.ChanSubscribe("$GVT.default.DP."+dataProduct+".*.EVENT.>", ch)

	if err != nil {
		return fmt.Errorf("subscribe failed %s", err.Error())
	}
	time.Sleep(1 * time.Second)
	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("unsubscribe failed %s", err.Error())
	}

	msg := <-ch

	err = proto.Unmarshal(msg.Data, &pe)
	if err != nil {
		fmt.Printf("Failed to parsing product event: %v", err)
	}

	r, err := pe.GetContent()
	if err != nil {
		fmt.Printf("Failed to parsing content: %v", err)
	}
	// fmt.Println(pe.EventName)
	// fmt.Println(r.AsMap())

	JSONByte, err := json.Marshal(r.AsMap())
	if err != nil {
		return fmt.Errorf("Receive payload marshal fail %s", err.Error())
	}
	recieveJSONStringResult := strings.Join(strings.Fields(string(JSONByte)), "")
	regexPayload := regexp.MustCompile(`'?([^']*)'?`).FindStringSubmatch(payload)[1]
	regexPayload = ut.FormatJSONData(regexPayload)

	if pe.EventName != event {
		return errors.New("event 比對不一致")
	}

	if recieveJSONStringResult != regexPayload {
		return errors.New("payload 資料不正確")
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
	ctx.Given(`^創建 data product "'(.*?)'" 使用參數 "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^"'(.*?)'" 創建 ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	ctx.When(`^publish Event "'(.*?)'" 使用參數 "'(.*?)'"$`, PublishEventCommand)
	ctx.Then(`^查詢 GVT_default_DP_"'(.*?)'" 裡有 "'(.*?)'" 帶有 "'(.*?)'"$`, CheckDPStreamDPExist)
	ctx.Then(`^使用 nats jetstream 查詢 GVT_default "'(.*?)'" 帶有 "'(.*?)'"$`, QueryJetstreamEventExist)
	ctx.When(`^更新 data product "'([^'"]*?)'" 使用參數 enabled=true$`, UpdateDataProductCommand)
	ctx.When(`^更新 data product "'([^'"]*?)'" 的 ruleset "'([^'"]*?)'" 使用參數 enabled=true$`, UpdateRulesetCommand)
	ctx.Then(`^查詢 GVT_default_DP_"'(.*?)'" 裡沒有 "'(.*?)'"$`, CheckDPStreamDPNotExist)
}
