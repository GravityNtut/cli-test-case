package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

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

func (testUtils *TestUtils) LoadConfig() error {
	str, err := os.ReadFile("../config/config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), &testUtils.Config)
	if err != nil {
		return err
	}
	return nil
}

func (testUtils *TestUtils) ExecuteShell(command string) error {
	f, err := os.Create("command.sh")
	f.WriteString(command)
	defer f.Close()

	cmd := exec.Command("sh", "./command.sh")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdResultTmp := &testUtils.CmdResult
	cmdResultTmp.Err = cmd.Run()

	cmdResultTmp.Stdout = stdout.String()
	cmdResultTmp.Stderr = stderr.String()
	return err
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

func (testUtils *TestUtils) ClearDataProducts() {
	nc, _ := nats.Connect("nats://" + testUtils.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	streams := js.StreamNames()

	re := regexp.MustCompile(`^GVT_default_DP_(.*)`)
	for stringName := range streams {
		parts := re.FindStringSubmatch(stringName)
		if parts == nil {
			continue
		}
		productName := parts[1]
		cmd := exec.Command("../gravity-cli", "product", "delete", productName, "-s", testUtils.Config.JetstreamURL)
		cmd.Run()
	}
}

func (testUtils *TestUtils) RestartDocker() {
	cmd := exec.Command("docker", "compose", "down")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err, stderr.String())
	}
	cmd = exec.Command("docker", "compose", "-f", "../docker-compose.yaml", "up", "-d")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err, stderr.String())
	}
}

func (testUtils *TestUtils) CheckNatsService() error {
	nc, err := nats.Connect("nats://" + testUtils.Config.JetstreamURL)
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

func (testUtils *TestUtils) CreateDataProduct(dataProduct string) error {
	cmd := exec.Command("../gravity-cli", "product", "create", dataProduct, "-s", testUtils.Config.JetstreamURL)
	return cmd.Run()
}

func (testUtils *TestUtils) CreateDataProductRuleset(dataProduct string, ruleset string) error {
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "add", dataProduct, ruleset, "--event", "test", "--method", "create", "-s", testUtils.Config.JetstreamURL)
	return cmd.Run()
}
