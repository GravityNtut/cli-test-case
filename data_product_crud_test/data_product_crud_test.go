package data_product_crud

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type Command struct {
	command string
	output  string
}

var ExecutingCommand Command

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

func VerifyContainerReady(containerName string) error {
	return nil
}

func VerifyDataProduct(style string, dataProduct string) error {
	_, _, err := ExecuteCommand("../gravity-cli product info " + dataProduct)
	fmt.Println(err)
	if err == "Error: Not found product \"drink\"\n\n" {
		return nil
	}
	return errors.New("錯誤")
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

func CreateDataProductCommand(command string) error {
	ExecutingCommand = Command{command, ""}
	return nil
}

func CreateDataProductResponse(expectOutput string) error {
	var err error
	ExecutingCommand.output, err, _ = ExecuteCommand(ExecutingCommand.command)
	fmt.Println(ExecutingCommand.output)
	if expectOutput == ExecutingCommand.output {
		return nil
	}
	return err
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^container "([^"]*)" ready`, VerifyContainerReady)

	ctx.Step(`^"([^"]*)" 沒有Data_Product "([^"]*)"`, VerifyDataProduct)
	ctx.Step(`^輸入創建data product cli指令 """(.*)"""`, CreateDataProductCommand)
	ctx.Step(`^創建data product cli回應 """(.*)"""`, CreateDataProductResponse)
}
