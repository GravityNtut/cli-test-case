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
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
)

type Config struct {
	JetstreamURL string
}

var config Config = Config{JetstreamURL: "172.25.19.180:32803"}

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

func CreateSchema(arg *godog.DocString) error {
	f, err := os.Create("schema.json")
	f.WriteString(arg.Content)
	defer f.Close()
	return err
}

func CreateDataProductCommand(dataProduct string, description string, state string) error {
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "--desc", description, "--enabled", "--schema", "schema.json", "-s", config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := stdout.String()

	fmt.Println("stdout: ", stdout.String())
	fmt.Println("stderr: ", stderr.String())

	if strings.EqualFold(state, "success") {
		if outStr == "Product \""+dataProduct+"\" was created\n" {
			return nil
		}
		return errors.New("應該要創建成功")
	} else if strings.EqualFold(state, "fail") {
		if err != nil {
			return nil
		}
		return errors.New("應該要創建失敗")
	}
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

	ctx.Given(`^schema=$`, CreateSchema)
	ctx.When(`^創建data product "([^"]*)" 註解 "([^"]*)" "([^"]*)"$`, CreateDataProductCommand)
	ctx.Then(`^使用gravity-cli 查詢 "([^"]*)" 成功$`, SearchDataProductByCLISuccess)
	ctx.Then(`^使用nats jetstream 查詢 "([^"]*)" 成功$`, SearchDataProductByJetstreamSuccess)
}
