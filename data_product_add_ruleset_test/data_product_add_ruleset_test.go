package data_product_ruleset_add

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
	cmd.Run()
	return nil
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

func AddRulesetCommand(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string) error {
	dataProduct = ProcessString(dataProduct)
	ruleset = ProcessString(ruleset)
	commandString := "../gravity-cli product ruleset add " + dataProduct + " " + ruleset
	if event != "[ignore]" {
		event := ProcessString(event)
		commandString += " --event " + event
	}
	if method != "[ignore]" {
		method := ProcessString(method)
		commandString += " --method " + method
	}
	if pk != "[ignore]" {
		pk := ProcessString(pk)
		commandString += " --pk " + pk
	}
	if desc != "[ignore]" {
		if desc == "[null]" {
			commandString += " --desc "
		} else {
			desc := ProcessString(desc)
			commandString += " --desc \"" + desc + "\""
		}
	}
	if handler != "[ignore]" {
		commandString += " --handler ./assets/" + handler
	}
	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " -s " + config.JetstreamURL
	ExecuteShell(commandString)
	return nil
}

func AddRulesetCommandFailed() error {
	if cmdResult.err != nil {
		return nil
	}
	return fmt.Errorf("ruleset 創建應該要失敗")
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

func AddRulesetCommandSuccess() error {
	if cmdResult.err == nil {
		return nil
	}
	return fmt.Errorf("ruleset 創建應該要成功")
}

func SearchRulesetByCLISuccess(dataProduct string, ruleset string) error {
	dataProduct = ProcessString(dataProduct)
	ruleset = ProcessString(ruleset)
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
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
	defer nc.Close()
	return err
}

func checkDispatcherService() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^已開啟服務nats$`, CheckNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, checkDispatcherService)

	ctx.Given(`^已有data product "([^"]*)"$`, CreateDataProduct)

	ctx.When(`^"([^"]*)" 創建ruleset "([^"]*)" method "([^"]*)" event "([^"]*)" pk "([^"]*)" desc "([^"]*)" handler "([^"]*)" schema "([^"]*)"$`, AddRulesetCommand)
	ctx.Then(`^ruleset 創建失敗$`, AddRulesetCommandFailed)

	ctx.Then(`^ruleset 創建成功$`, AddRulesetCommandSuccess)
	ctx.Then(`^使用gravity-cli 查詢 "([^"]*)" 的 "([^"]*)" 成功$`, SearchRulesetByCLISuccess)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
}
