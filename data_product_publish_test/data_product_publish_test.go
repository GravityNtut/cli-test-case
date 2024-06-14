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

func CreateDataProductRuleset(dataProduct string, ruleset string, RSMethod string, RSPk string, RSHandler string, RSSchema string, RSEnabled string) error {

	if RSEnabled == testutils.TrueString {
		RSEnabled = "true"
	} else if RSEnabled == testutils.FalseString {
		RSEnabled = "false"
	} else {
		return errors.New("Enable must be [true] or [false]")
	}

	if RSHandler == testutils.IgnoreString {
		RSHandler = ""
	} else {
		RSHandler = "--handler=" + RSHandler
	}

	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", ruleset, "--method", RSMethod, "--pk", RSPk, RSHandler, "--schema", RSSchema, "--enabled="+RSEnabled, "-s", ut.Config.JetstreamURL)
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

	var m *nats.Msg
	select {
	case m = <-ch:

	case <-time.After(10 * time.Second):
		return errors.New("subscribe out of time")
	}

	var jsonData JSONData
	if err := json.Unmarshal(m.Data, &jsonData); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}
	//fmt.Println("Event: ", jsonData.Event)
	result, _ := Base64ToString(jsonData.Payload)
	//fmt.Println("Payload: ", result)

	if jsonData.Event != event {
		return fmt.Errorf("expected event: %s, actual event: %s", event, jsonData.Event)
	}
	if result != payload {
		return fmt.Errorf("expected payload: %s, actual payload: %s", payload, result)
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

	select {
	case <-ch:
		return fmt.Errorf("expected not publish in GVT_default_DP，but now in GVT_default_DP")
	case <-time.After(5 * time.Second):

	}
	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("unsubscribe failed %s", err.Error())
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
	//fmt.Println(stderr.String())
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
	//fmt.Println(stderr.String())
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

	var msg *nats.Msg

	select {
	case msg = <-ch:

	case <-time.After(10 * time.Second):
		return errors.New("subscribe out of time")
	}

	err = proto.Unmarshal(msg.Data, &pe)
	if err != nil {
		return fmt.Errorf("Failed to parsing product event: %v", err)
	}

	r, err := pe.GetContent()
	if err != nil {
		return fmt.Errorf("Failed to parsing content: %v", err)
	}

	JSONByte, err := json.Marshal(r.AsMap())
	if err != nil {
		return fmt.Errorf("Receive payload marshal fail %s", err.Error())
	}
	recieveJSONStringResult := strings.Join(strings.Fields(string(JSONByte)), "")
	regexPayload := regexp.MustCompile(`'?([^']*)'?`).FindStringSubmatch(payload)[1]

	regexPayload = ut.FormatJSONData(regexPayload)

	if pe.EventName != event {
		return errors.New("event is not matched")
	}

	var receivedMap map[string]interface{}
	var payloadMap map[string]interface{}

	err = json.Unmarshal([]byte(recieveJSONStringResult), &receivedMap)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal received JSON: %s", err.Error())
	}

	err = json.Unmarshal([]byte(regexPayload), &payloadMap)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal payload JSON: %s", err.Error())
	}

	filteredMap := filterMapByKeys(receivedMap, payloadMap)

	filteredJSON, err := json.Marshal(filteredMap)
	if err != nil {
		return fmt.Errorf("Failed to marshal filtered JSON: %s", err.Error())
	}

	filteredJSONStringResult := strings.Join(strings.Fields(string(filteredJSON)), "")

	// fmt.Println("recieveJSONStringResult:", recieveJSONStringResult)
	// fmt.Println("filteredJSONStringResult: ",filteredJSONStringResult)

	if filteredJSONStringResult != recieveJSONStringResult {
		return errors.New("payload is not matched")
	}

	return nil
}

func filterMapByKeys(source, keys map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})
	for key := range keys {
		if value, exists := source[key]; exists {
			filtered[key] = value
		}
	}
	return filtered
}

func CheckDPStreamDPEventHasTwoPayload(dataProduct string, event string, payload string, payload2 string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	var pe gravity_sdk_types_product_event.ProductEvent

	var payloadList = []string{payload, payload2}

	channel := make(chan error)

	ch := make(chan *nats.Msg, 2)
	sub, err := js.ChanSubscribe("$GVT.default.DP."+dataProduct+".*.EVENT.>", ch)
	if err != nil {
		return fmt.Errorf("js get subscribe failed: %v", err)
	}
	time.Sleep(1 * time.Second)
	if err := sub.Unsubscribe(); err != nil {
		return fmt.Errorf("unsubscribe failed %s", err.Error())
	}

	go func() {
		i := 0
		for msg := range ch {
			err = proto.Unmarshal(msg.Data, &pe)
			if err != nil {
				channel <- fmt.Errorf("Failed to parsing product event: %v", err)
			}

			r, err := pe.GetContent()
			if err != nil {
				channel <- fmt.Errorf("Failed to parsing content: %v", err)
			}
			JSONByte, err := json.Marshal(r.AsMap())
			if err != nil {
				channel <- fmt.Errorf("Receive payload marshal fail %s", err.Error())
			}
			recieveJSONStringResult := strings.Join(strings.Fields(string(JSONByte)), "")
			regexPayload := regexp.MustCompile(`'?([^']*)'?`).FindStringSubmatch(payloadList[i])[1]
			regexPayload = ut.FormatJSONData(regexPayload)
			if recieveJSONStringResult != regexPayload {
				channel <- fmt.Errorf("payload is not matched")
			}
			i++
		}
		channel <- nil
	}()

	time.Sleep(10 * time.Second)

	close(ch)

	if err := <-channel; err != nil {
		return err
	}

	if pe.EventName != event {
		return errors.New("event is not matched")
	}
	return nil
}

func PublishEventCommandFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return fmt.Errorf("publish should be failed")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^Nats has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product with "'(.*?)'" using parameters "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^"'(.*?)'" create ruleset "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	ctx.When(`^Publish Event "'(.*?)'" using parameters "'(.*?)'"$`, PublishEventCommand)
	ctx.Then(`^Query GVT_default_DP_"'(.*?)'" has "'(.*?)'" with "'(.*?)'"$`, CheckDPStreamDPExist)
	ctx.Then(`^Using NATS Jetstream to query GVT_default "'(.*?)'" has "'(.*?)'"$`, QueryJetstreamEventExist)
	ctx.When(`^Update data product "'([^'"]*?)'" using parameters enabled=true$`, UpdateDataProductCommand)
	ctx.When(`^Update data product "'([^'"]*?)'" ruleset "'([^'"]*?)'" using parameters enabled=true$`, UpdateRulesetCommand)
	ctx.Then(`^Query GVT_default_DP_"'(.*?)'" has no "'(.*?)'"$`, CheckDPStreamDPNotExist)
	ctx.Then(`^Query GVT_default_DP_"'(.*?)'" has "'(.*?)'" with "'(.*?)'" and "'(.*?)'"$`, CheckDPStreamDPEventHasTwoPayload)
	ctx.Then(`^Cli returns create failed$`, PublishEventCommandFailed)
}
