package data_product_delete

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
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string
}

type CmdResult struct {
	err    error
	stdout string
	stderr string
}

var config Config = Config{JetstreamURL: "0.0.0.0:32803"}
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
	// LoadConfig()
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

func SearchDataProductByCLIFail(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil
	}
	return errors.New("Cli 預期不存在，但查詢到了")
}

func SearchDataProductByJetstreamFail(dataProduct string) error {
	nc, _ := nats.Connect("nats://" + config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return errors.New("Jetstream 預期不存在，但查詢到了")
		}
	}
	return nil
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

func CreateDataProduct(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "-s", config.JetstreamURL)
	cmd.Run()
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
	return nil
}

func DeleteDataProductCommand(dataProduct string) error {
	commandString := "../gravity-cli product delete " + dataProduct + " -s " + config.JetstreamURL
	ExecuteShell(commandString)
	return nil
}

func DeleteDataProductSuccess() error {
	return cmdResult.err
}

func DeleteDataProductFail() error {
	if cmdResult.err != nil {
		return nil
	}
	return errors.New("data product 刪除應該要失敗")
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^已開啟服務nats$`, checkNatsService)
	ctx.Given(`^已開啟服務dispatcher$`, checkDispatcherService)

	ctx.Given(`^已有data product "([^"]*)"$`, CreateDataProduct)
	ctx.When(`^刪除data product "([^"]*)"$`, DeleteDataProductCommand)
	ctx.Then(`^data product 刪除成功$`, DeleteDataProductSuccess)
	ctx.Then(`^data product 刪除失敗$`, DeleteDataProductFail)
	ctx.Then(`^使用gravity-cli查詢data product 列表 "([^"]*)" 不存在$`, SearchDataProductByCLIFail)
	ctx.Then(`^使用nats jetstream 查詢 data product 列表 "([^"]*)" 不存在$`, SearchDataProductByJetstreamFail)
}
