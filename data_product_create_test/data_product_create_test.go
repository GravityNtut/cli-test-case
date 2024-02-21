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
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string
}

type CreateProductCommand struct {
	name   string
	desc   string
	schema string
}

var config Config = Config{JetstreamURL: "0.0.0.0:32803"}
var createProductCommand CreateProductCommand

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
			StopOnFailure: true,
			TestingT:      t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func CreateDataProductCommand(dataProduct string, description string, schema string) error {
	createProductCommand.name = dataProduct
	createProductCommand.desc = description
	createProductCommand.schema = schema
	return nil
}

func ExecucteCreateProductCommandSuccess() error {
	dataProduct := createProductCommand.name
	description := createProductCommand.desc
	schema := createProductCommand.schema

	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "--desc", description, "--enabled", "--schema", schema, "-s", config.JetstreamURL)
	err := cmd.Run()

	return err
}

func SearchDataProductByCLISuccess(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct, "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
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

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})
	ctx.Step(`^創建一個data product "([^"]*)" 註解 "([^"]*)" schema檔案 "([^"]*)"$`, CreateDataProductCommand)
	ctx.Step("^建立成功$", ExecucteCreateProductCommandSuccess)
	ctx.Step(`^使用gravity-cli查詢data product 列表 "([^"]*)" 存在$`, SearchDataProductByCLISuccess)
	ctx.Step(`^使用nats jetstream 查詢 data product 列表 "([^"]*)" 存在$`, SearchDataProductByJetstreamSuccess)

}
