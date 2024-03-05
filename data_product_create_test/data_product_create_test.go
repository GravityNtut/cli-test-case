package data_product_create

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string `json:"jetstream_url"`
}

type CmdResult struct {
	stdout string
	stderr string
}

var config Config
var cmdResult CmdResult

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

func ExecuteShell(command string) error {
	f, err := os.Create("command.sh")
	f.WriteString(command)
	defer f.Close()

	cmd := exec.Command("sh", "./command.sh")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Run()

	cmdResult.stdout = stdout.String()
	cmdResult.stderr = stderr.String()
	return err
}

func CreateDataProductCommand(dataProduct string, description string, schema string) error {
	dataProduct = ProcessString(dataProduct)
	commandString := "../gravity-cli product create "
	if dataProduct != "[null]" {
		commandString += dataProduct
	}
	if description != "[ignore]" {
		if description == "[null]" {
			commandString += " --desc"
		} else {
			description := ProcessString(description)
			commandString += " --desc \"" + description + "\""
		}
	}

	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " --enabled" + " -s " + config.JetstreamURL
	ExecuteShell(commandString)

	return nil
}

func CreateDateProductCommandSuccess(productName string) error {
	outStr := cmdResult.stdout
	productName = ProcessString(productName)
	if outStr == "Product \""+productName+"\" was created\n" {
		return nil
	}
	return errors.New("Cli回傳訊息錯誤")
}

func CreateDateProductCommandFail() error {
	outErr := cmdResult.stderr
	outStr := cmdResult.stdout
	if outStr == "" && outErr != "" {
		return nil
	}
	return errors.New("Cli回傳訊息錯誤")
}

func SearchDataProductByCLISuccess(dataProduct string) error {
	dataProduct = ProcessString(dataProduct)
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err
}

func ProcessString(s string) string {

	re := regexp.MustCompile(`\[(\S+)\]x(\d+)`)
	parts := re.FindStringSubmatch(s)
	if parts == nil {
		return s
	}
	str := parts[1]
	times, _ := strconv.Atoi(parts[2])
	completeString := ""
	for i := 0; i < times; i++ {
		completeString += str
	}
	return completeString
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	dataProduct = ProcessString(dataProduct)
	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return nil
		}
	}
	return errors.New("jetstream裡未創建成功")
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

func AssertErrorMessages(errorMessage string) error {
	// TODO
	// outErr := cmdResult.stderr
	// if outErr == errorMessage {
	// 	return nil
	// }
	// return errors.New("Cli回傳訊息錯誤")
	return nil
}

func checkNatsService() error {
	nc, err := nats.Connect("nats://" + config.JetstreamURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	return nil
}

func checkDispatcherService() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println(container.Names[0])
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

	ctx.Given(`^已開啟服務nats$`, checkNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, checkDispatcherService)
	ctx.When(`^創建一個data product "([^"]*)" 註解 "([^"]*)" schema檔案 "([^"]*)"$`, CreateDataProductCommand)
	ctx.Step(`^Cli回傳"([^"]*)"建立成功$`, CreateDateProductCommandSuccess)
	ctx.Step(`^Cli回傳建立失敗$`, CreateDateProductCommandFail)
	ctx.Step(`^使用gravity-cli查詢data product 列表 "([^"]*)" 存在$`, SearchDataProductByCLISuccess)
	ctx.Step(`^使用nats jetstream 查詢 data product 列表 "([^"]*)" 存在$`, SearchDataProductByJetstreamSuccess)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
}
