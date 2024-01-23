package data_product_create

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cucumber/godog"
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

func ExecuteCommand(command string) (string, error, string) {
	fmt.Println(command)
	comArr := strings.Fields(command)
	com := comArr[0]
	args := comArr[1:]
	cmd := exec.Command(com, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := stdout.String()
	return outStr, err, stderr.String()
}

func ClearDataProduct() error {
	return nil
}

func CreateSchema(arg *godog.DocString) error {
	fmt.Println("here")
	fmt.Println(arg.Content)
	f, err := os.Create("schema.json")
	f.WriteString(arg.Content)
	defer f.Close()
	return err
}

func CreateDataProductCommand(dataProduct string, description string) error {
	exec.Command("../gravity-cli", "product", "create", dataProduct, "--desc", description, "--enabled")
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	//ctx.Step(`沒有data product`, ClearDataProduct)

	ctx.Step(`schema.json=`, CreateSchema)
	ctx.Step(`創建data product "([^"]*)" 註解 "([^"]*)"`, CreateDataProductCommand)
}
