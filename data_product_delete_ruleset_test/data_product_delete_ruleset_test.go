package data_product_ruleset_delete

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/cucumber/godog"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string `json:"jetstream_url"`
}

type CommandResult struct {
	err    error
	Stdout string
	Stderr string
}

var config Config
var cmdResult CommandResult

func LoadConfig() error {
	str, err := os.ReadFile("../config/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), &config)
	if err != nil {
		return err
	}

	return nil
}
func TestFeatures(t *testing.T) {
	LoadConfig()
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

func AddRulesetCommand(dataProduct string, ruleset string) error {
	commandString := "../gravity-cli product ruleset add " + dataProduct + " " + ruleset + " --event drinkCreated --method create --pk id --desc \"desc\" --handler ./assets/handler.js --schema ./assets/schema.json -s " + config.JetstreamURL
	ExecuteShell(commandString)
	return nil
}

func DeleteRulesetCommand(productName string, rulesetName string) error {
	commandString := "../gravity-cli product ruleset delete " + productName + " " + rulesetName + " -s " + config.JetstreamURL
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
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil
	}
	return fmt.Errorf("ruleset 應該不存在")
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.Names[0] == "/gravity-dispatcher" {
			return nil
		}
	}
	return errors.New("dispatcher container 不存在")
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
}