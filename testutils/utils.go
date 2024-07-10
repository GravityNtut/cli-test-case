package testutils

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
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

type TestUtils struct {
	Config    Config
	CmdResult CommandResult
}

type CommandResult struct {
	Err    error
	Stdout string
	Stderr string
}

type Config struct {
	JetstreamURL  string `json:"jetstream_url"`
	StopOnFailure bool   `json:"stop_on_failure"`
}

const (
	NullString       = "[null]"
	IgnoreString     = "[ignore]"
	TrueString       = "[true]"
	FalseString      = "[false]"
	NatsProtocol     = "nats://"
	GravityCliString = "../gravity-cli"
)

func (testUtils *TestUtils) LoadConfig() error {
	str, err := os.ReadFile("../config/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, &testUtils.Config)
	if err != nil {
		return err
	}
	return nil
}

func (testUtils *TestUtils) ExecuteShell(command string) error {
	f, err := os.Create("command.sh")
	if err != nil {
		return err
	}
	_, err = f.WriteString(command)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	cmd := exec.Command("sh", "./command.sh")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdResultTmp := &testUtils.CmdResult
	cmdResultTmp.Err = cmd.Run()

	cmdResultTmp.Stdout = stdout.String()
	cmdResultTmp.Stderr = stderr.String()
	return nil
}

func (testUtils *TestUtils) ProcessString(str string) string {
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

func (testUtils *TestUtils) AssertStringEqual(actual string, expected string) error {
	if expected != actual {
		return fmt.Errorf("expect: %s actual: %s", expected, actual)
	}
	return nil
}

func (testUtils *TestUtils) AssertIntEqual(actual int, expected int) error {
	if expected != actual {
		return fmt.Errorf("expect: %d actual: %d", expected, actual)
	}
	return nil
}

func (testUtils *TestUtils) AssertUIntEqual(actual uint64, expected uint64) error {
	if expected != actual {
		return fmt.Errorf("expect: %d actual: %d", expected, actual)
	}
	return nil
}

func (testUtils *TestUtils) ValidateField(actual, expected string) error {
	if expected != IgnoreString {
		regex := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(expected)[1] //移除雙引號
		if actual != regex {
			return fmt.Errorf("%s 與nats資訊不符", expected)
		}
	}
	return nil
}

func (testUtils *TestUtils) ValidateEnabled(actual bool, expected string) error {
	var enabledBool bool
	if expected == TrueString {
		enabledBool = true
	} else if expected == IgnoreString || expected == FalseString {
		enabledBool = false
	} else {
		return errors.New("不允許true,false,ignore以外的輸入")
	}
	if enabledBool != actual {
		return errors.New("enabled 更改失敗")

	}
	return nil
}

func (testUtils *TestUtils) ValidateHandler(actual interface{}, expected string) error {
	if expected != IgnoreString {
		regexHandler := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(expected)[1]
		fileContent, err := os.ReadFile(regexHandler)
		if err != nil {
			return err
		}
		rulesetHandler := actual.(map[string]interface{})
		handlerScript := rulesetHandler["script"].(string)
		if string(fileContent) != handlerScript {
			return errors.New("handler與nats資訊不符")
		}
	}
	return nil
}

func (testUtils *TestUtils) ValidateSchema(actual interface{}, expected string) error {
	if expected != IgnoreString {
		regexSchema := regexp.MustCompile(`"?([^"]*)"?`).FindStringSubmatch(expected)[1]
		fileContent, err := os.ReadFile(regexSchema)
		if err != nil {
			return err
		}
		natsSchema, _ := json.Marshal(actual)
		fileSchema := testUtils.FormatJSONData(string(fileContent))
		if fileSchema != string(natsSchema) {
			return errors.New("schema與nats資訊")
		}
	}
	return nil
}

func (testUtils *TestUtils) FormatJSONData(JSONData string) string {
	var JSON interface{}
	err := json.Unmarshal([]byte(JSONData), &JSON)
	if err != nil {
		log.Fatalf("%s Unmarshal Fail %s", JSONData, err.Error())
	}
	JSONByte, err := json.Marshal(JSON)
	if err != nil {
		log.Fatalf("%s Marshal Fail %s", JSONData, err.Error())
	}
	expectJSONStringResult := strings.Join(strings.Fields(string(JSONByte)), "")
	return expectJSONStringResult
}

func (testUtils *TestUtils) ClearDataProducts() {
	nc, _ := nats.Connect(NatsProtocol + testUtils.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames(nats.MaxWait(200 * time.Second))

	re := regexp.MustCompile(`^GVT_default_DP_(.*)`)
	for stringName := range streams {
		if stringName == "GVT_default" {
			if err := js.PurgeStream(stringName); err != nil {
				log.Fatalf(err.Error())
			}
		}
		parts := re.FindStringSubmatch(stringName)
		if parts == nil {
			continue
		}
		productName := parts[1]
		cmd := exec.Command(GravityCliString, "product", "delete", productName, "-s", testUtils.Config.JetstreamURL)
		if err := cmd.Run(); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func (testUtils *TestUtils) RestartDocker() {
	cmd := exec.Command("docker", "compose", "down")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), stderr.String())
	}
	cmd = exec.Command("docker", "compose", "-f", "../docker-compose.yaml", "up", "-d")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), stderr.String())
	}
}

func (testUtils *TestUtils) CheckNatsService() error {
	nc, err := nats.Connect(NatsProtocol + testUtils.Config.JetstreamURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	return nil
}

func (testUtils *TestUtils) CheckDispatcherService() error {
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

func (testUtils *TestUtils) CreateDataProduct(dataProduct string, enabled string) error {
	if enabled == TrueString {
		cmd := exec.Command(GravityCliString, "product", "create", dataProduct, "--enabled", "-s", testUtils.Config.JetstreamURL)
		return cmd.Run()
	} else if enabled == IgnoreString || enabled == FalseString {
		cmd := exec.Command(GravityCliString, "product", "create", dataProduct, "-s", testUtils.Config.JetstreamURL)
		return cmd.Run()
	}
	return errors.New("dataProduct create 參數錯誤")
}

func (testUtils *TestUtils) CreateDataProductRuleset(dataProduct string, ruleset string, enabled string) error {
	if enabled == TrueString {
		cmd := exec.Command(GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--enabled", "--event", "test", "--method", "create", "-s", testUtils.Config.JetstreamURL)
		return cmd.Run()
	} else if enabled == IgnoreString || enabled == FalseString {
		cmd := exec.Command(GravityCliString, "product", "ruleset", "add", dataProduct, ruleset, "--event", "test", "--method", "create", "-s", testUtils.Config.JetstreamURL)
		return cmd.Run()
	}
	return errors.New("ruleset add 參數錯誤")
}
