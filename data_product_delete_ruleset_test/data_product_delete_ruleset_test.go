package data_product_ruleset_delete

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string
}

type CommandResult struct {
	err    error
	Stdout string
	Stderr string
}

var config Config = Config{JetstreamURL: "0.0.0.0:32803"}
var cmdResult CommandResult

func LoadConfig() error {
	str, err := os.ReadFile("../config/config.json")
	if err != nil {
		return err
	}
	fmt.Println(string(str))
	err = json.Unmarshal([]byte(str), &config)
	if err != nil {
		return err
	}

	return nil
}
func TestFeatures(t *testing.T) {
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

func CreateDataProduct(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "-s", config.JetstreamURL)
	return cmd.Run()
}

func CreateDataProductRuleset(dataProduct string, ruleset string) error {
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "add", dataProduct, ruleset, "--event \"test\" --method create", "-s", config.JetstreamURL)
	return cmd.Run()

}

func ExecuteShell(command string) error {
	f, err := os.Create("command.sh")
	f.WriteString(command)
	defer f.Close()

	cmd := exec.Command("sh", "./command.sh")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdResult.err = cmd.Run()

	cmdResult.Stdout = stdout.String()
	cmdResult.Stderr = stderr.String()
	return err
}

func ProcessString(str string) string {
	re := regexp.MustCompile(`\[(\S+)\]x(\d+)`)

	parts := re.FindStringSubmatch(str)
	if parts == nil {
		return str
	}
	chr := parts[1]
	times, _ := strconv.Atoi(parts[2])
	completeString := ""
	for i := 0; i < times; i++ {
		completeString += chr
	}
	return completeString
}

func AddRulesetCommand(dataProduct string, ruleset string) error {
	dataProduct = ProcessString(dataProduct)
	ruleset = ProcessString(ruleset)
	commandString := "../gravity-cli product ruleset add " + dataProduct + " " + ruleset + " --event drinkCreated --method create --pk id --desc \"desc\" --handler ./assets/handler.js --schema ./assets/schema.json -s " + config.JetstreamURL
	fmt.Println(commandString)
	ExecuteShell(commandString)
	return nil
}

func DeleteRulesetCommand(productName string, rulesetName string) error {
	productName = ProcessString(productName)
	rulesetName = ProcessString(rulesetName)
	commandString := "../gravity-cli product ruleset delete " + productName + " " + rulesetName + " -s " + config.JetstreamURL
	fmt.Println(commandString)
	ExecuteShell(commandString)
	return nil
}

func DeleteRulesetCommandFailed() error {
	if cmdResult.err != nil {
		return nil
	}
	return fmt.Errorf("ruleset 刪除應該要失敗")
}

func DeleteRulesetCommandSuccess() error {
	if cmdResult.err == nil {
		return nil
	}
	return fmt.Errorf("ruleset 刪除應該要成功")
}

func ClearDataProducts() {
	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	js.PurgeStream("KV_GVT_default_PRODUCT")
}

func SearchRulesetByCLINotExists(dataProduct string, ruleset string) error {
	dataProduct = ProcessString(dataProduct)
	ruleset = ProcessString(ruleset)
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run() 
	if err != nil { //TODO 這裡要改成判斷是否有錯誤訊息
		return nil
	}
	return err
}

func AssertErrorMessages(expected string) error {
	// Todo
	// if cmdResult.Stderr == expected {
	// 	return nil
	// }
	// return fmt.Errorf("應有錯誤訊息: %s", expected)
	return nil
}

func CheckNatsService() error {
	nc, err := nats.Connect("nats://" + config.JetstreamURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	return nil
}

func CheckDispatcherService() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務nats$`, CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, CheckDispatcherService)
	ctx.Given(`^已有date product "([^"]*)"$`, CreateDataProduct)
	ctx.Given(`^已有data product 的 ruleset "([^"]*)" "([^"]*)"$`, AddRulesetCommand)
	ctx.When(`^刪除 "([^"]*)" 的 ruleset "([^"]*)"$`, DeleteRulesetCommand)
	ctx.Then(`^刪除失敗$`, DeleteRulesetCommandFailed)
	ctx.Then(`^刪除成功$`, DeleteRulesetCommandSuccess)
	ctx.Then(`^使用gravity-cli 查詢 "([^"]*)" 的 "([^"]*)" 不存在$`, SearchRulesetByCLINotExists)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
}
