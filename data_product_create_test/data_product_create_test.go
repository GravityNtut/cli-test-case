package data_product_create

import (
	"bytes"
	"context"
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

func TestFeatures(t *testing.T) {
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
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "--desc", description, "--enabled", "--schema", "schema.json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := stdout.String()

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
	cmd := exec.Command("../gravity-cli", "product", "info", dataProduct)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return err
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	nc, _ := nats.Connect("nats://127.0.0.1:32803")
	defer nc.Drain()

	// Request the list of all streams
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	// Fetch the list of streams
	streams := js.StreamNames()

	for stream := range streams {
		if stream == "GVT_default_DP_"+dataProduct {
			return nil
		}
	}
	return errors.New("jetstream裡未創建成功")
}

func ClearDataProducts() {
	nc, _ := nats.Connect("nats://127.0.0.1:32803")
	defer nc.Drain()

	// Request the list of all streams
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	for kv := range js.KeyValueStoreNames() {
		fmt.Println(kv)
	}

	for stream := range js.StreamNames() {
		js.DeleteStream(stream)
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	// 	cmd := exec.Command("docker", "compose", "-f", "./docker-compose.yaml", "up", "-d")
	// 	if err := cmd.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return ctx, nil
	// })

	// ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	// 	cmd := exec.Command("docker", "compose", "down")
	// 	if err := cmd.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return ctx, nil
	// })

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ClearDataProducts()
		return ctx, nil
	})

	ctx.Step(`^schema=$`, CreateSchema)
	ctx.Step(`^創建data product "([^"]*)" 註解 "([^"]*)" "([^"]*)"$`, CreateDataProductCommand)
	ctx.Step(`^使用gravity-cli 查詢 "([^"]*)" 成功$`, SearchDataProductByCLISuccess)
	ctx.Step(`^使用nats jetstream 查詢 "([^"]*)" 成功$`, SearchDataProductByJetstreamSuccess)
}
