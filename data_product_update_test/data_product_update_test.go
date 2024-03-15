package data_product_update

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
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

type JsonData struct {
	Name    string      `json:"name"`
	Desc    string      `json:"desc"`
	Enabled bool        `json:"enabled"`
	Schema  interface{} `json:"schema"`
}

type Config struct {
	JetstreamURL string `json:"jetstream_url"`
}

type CmdResult struct {
	err    error
	stdout string
	stderr string
}

var originJsonData string
var newJsonData JsonData
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

	cmdResult.err = cmd.Run()

	cmdResult.stdout = stdout.String()
	cmdResult.stderr = stderr.String()

	return err
}

func UpdateDataProductCommand(dataProduct string, description string, enable string, schema string) error {
	dataProduct = ProcessString(dataProduct)
	commandString := "../gravity-cli product update "
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

	if enable != "[ignore]" {
		if enable == "[true]" {
			commandString += " --enabled"
		} else if enable == "" {
			commandString += ""
		} else {
			return errors.New("不允許true或ignore以外的輸入")
		}
	}

	if schema != "[ignore]" {
		commandString += " --schema ./assets/" + schema
	}
	commandString += " -s " + config.JetstreamURL
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

func CreateDataProduct(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "-s", config.JetstreamURL)
	cmd.Run()

	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	originJsonData = string(entry.Value())
	return nil
}

func DataProductNotChanges(dataProduct string, description string, schema string, enabled string) error {

	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	if string(entry.Value()) == originJsonData {
		return nil
	}
	return errors.New("與原始資料不符")
}

func DataProductUpdateSuccess(dataProduct string, description string, schema string, enabled string) error {
	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &newJsonData)
	if err != nil {
		fmt.Println("解碼 JSON 時出現錯誤:", err)
		return err
	}

	if schema != "[ignore]" {
		fileContent, err := os.ReadFile("./assets/" + schema)
		if err != nil {
			return err
		}
		schemaString, _ := json.Marshal(newJsonData.Schema)
		fileSchema := strings.TrimSpace(string(fileContent))
		if fileSchema != string(schemaString) {
			fmt.Println("File " + fileSchema)
			fmt.Println("schema " + string(schemaString))
			return errors.New("schema內容不同")
		}
	}

	enabledBool := false
	if enabled == "[true]" {
		enabledBool = true
	}

	if dataProduct == newJsonData.Name && description == newJsonData.Desc && enabledBool == newJsonData.Enabled {
		return nil
	}
	return errors.New("資料更新失敗")
}

func UpdateSuccess() error {
	return cmdResult.err
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務nats$`, checkNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, checkDispatcherService)
	ctx.Given(`^已有data product "([^"]*)"$`, CreateDataProduct)
	ctx.When(`^更新data product "([^"]*)" 註解 "([^"]*)" "([^"]*)" schema檔案 "([^"]*)"$`, UpdateDataProductCommand)
	ctx.Then(`^data product更改成功$`, UpdateSuccess)
	ctx.Then(`^Cli回傳建立失敗$`, CreateDateProductCommandFail)
	ctx.Then(`^應有錯誤訊息 "([^"]*)"$`, AssertErrorMessages)
	ctx.Then(`^使用nats驗證data product "([^"]*)" description "([^"]*)" schema檔案 "([^"]*)" enabled "([^"]*)" 更改成功`, DataProductUpdateSuccess)
	ctx.Then(`^使用nats驗證data product "([^"]*)" description "([^"]*)" schema檔案 "([^"]*)" enabled "([^"]*)" 資料無任何改動$`, DataProductNotChanges)
}
