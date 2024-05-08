package dataproductpublish

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

func CreateDataProductRuleset(dataProduct string, ruleset string, RSMethod string, event string, RSPk string, RSHandler string, RSSchema string, RSEnabled string) error {

	if RSEnabled == testutils.TrueString {
		RSEnabled = "true"
	} else if RSEnabled == testutils.FalseString {
		RSEnabled = "false"
	} else {
		return errors.New("Enable 必須要[true] 或 [false]")
	}
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", event, "--method", RSMethod, "--pk", RSPk, "--handler", RSHandler, "--schema", RSSchema, "--enabled="+RSEnabled, "-s", ut.Config.JetstreamURL)
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
	}
	return nil
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
	payload = ut.ProcessString(payload)
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *nats.Msg, 9)
	js.ChanSubscribe("$GVT.default.EVENT.*", ch)
	for i := 0; i < 9; i++ {
		m := <-ch
		fmt.Println("Received Num: ", i)
		var jsonData JSONData
		err := json.Unmarshal(m.Data, &jsonData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Event: ", jsonData.Event)
		result, _ := Base64ToString(jsonData.Payload)
		fmt.Println("Payload: ", result)
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

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務 nats$`, ut.CheckNatsService)
	ctx.Given(`^已開啟服務 dispatcher$`, ut.CheckDispatcherService)
	ctx.Given(`^創建 data product "'(.*?)'" 使用參數 "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^"'(.*?)'" 創建 ruleset "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	// ctx.When(`^publish Event "'(.*?)'" 使用參數 "'(.*?)'"$`, PublishEventCommand)
	// ctx.Then(`^使用 SDK 查詢 GVT_default_DP_drink 裡沒有 "'(.*?)'" 帶有 "'(.*?)'"$`)
	ctx.Then(`^使用 nats jetstream 查詢 GVT_default "'(.*?)'" 帶有 "'(.*?)'"$`, QueryJetstreamEventExist)
}
