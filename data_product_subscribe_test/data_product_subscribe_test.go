package dataproductsubscribe

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
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

func SubscribeDataProductCommand(productName string, subName string, partitions string, seq string) error {
	productName = ut.ProcessString(productName)
	subName = ut.ProcessString(subName)
	cmdString := "timeout 1 ../gravity-cli product sub "
	if productName != testutils.NullString {
		cmdString += productName + " "
	}
	if subName != testutils.IgnoreString {
		cmdString += "--name " + subName + " "
	}
	if partitions != testutils.IgnoreString {
		cmdString += "--partitions " + partitions + " "
	}
	if seq != testutils.IgnoreString {
		cmdString += "--seq " + seq + " "
	}

	cmdString += "-s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(cmdString)
	if err != nil {
		return err
	}
	return nil
}

func DisplayData() error {
	fmt.Println(ut.CmdResult.Stdout)
	return nil
}

func PublishProductEvent(eventName string) error {
	for i := 0; i < 10; i++ {
		payload := `{"id":%d, "name":"test%d", "kcal":%d, "price":%d}`
		result := fmt.Sprintf(payload, i+1, i+1, i*100, i+20)
		fmt.Println(result)
		fmt.Println(testutils.GravityCliString, "pub", eventName, result, "-s", ut.Config.JetstreamURL)
		cmd := exec.Command(testutils.GravityCliString, "pub", eventName, result, "-s", ut.Config.JetstreamURL)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

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
	method := "create"
	event := ruleset
	pk := "id"
	rulesetDesc := "drink創建事件"
	handler := "./assets/handler.js"
	schema := "./assets/schema.json"
	enabled := "true"
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", event, "--method", method, "--desc", rulesetDesc, "--pk", pk, "--handler", handler, "--schema", schema, "--enabled="+enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
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
	ctx.Given(`^已有 data product "'(.*?)'"$`, CreateDataProduct)
	ctx.Given(`^已有 data product 的 ruleset "'(.*?)'" "'(.*?)'"$`, CreateDataProductRuleset)
	ctx.Given(`^已 publish 10 筆 "'(.*?)'" Event$`, PublishProductEvent)
	ctx.When(`^訂閱data product "'(.*?)'" 使用參數 "'(.*?)'" "'(.*?)'" "'(.*?)'"`, SubscribeDataProductCommand)
	ctx.Then(`^顯示資料`, DisplayData)
}
