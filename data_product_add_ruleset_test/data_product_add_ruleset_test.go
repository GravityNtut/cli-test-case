package dataproductrulesetadd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

type Rule struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Desc       string      `json:"desc"`
	Event      string      `json:"event"`
	Product    string      `json:"product"`
	Method     string      `json:"method"`
	PrimaryKey []string    `json:"primaryKey"`
	Enabled    bool        `json:"enabled"`
	Handler    interface{} `json:"handler"`
	Schema     interface{} `json:"schema"`
	CreatedAt  string      `json:"createdAt"`
	UpdatedAt  string      `json:"updatedAt"`
}

type RuleMap map[string]Rule

type JSONData struct {
	Name            string      `json:"name"`
	Desc            string      `json:"desc"`
	Enabled         bool        `json:"enabled"`
	Rules           RuleMap     `json:"rules"`
	Schema          interface{} `json:"schema"`
	EnabledSnapshot bool        `json:"enabledSnapshot"`
	Snapshot        interface{} `json:"snapshot"`
	Stream          string      `json:"stream"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
}

var jsonData JSONData

var ut = testutils.TestUtils{}

var opts = godog.Options{
	Format:        "pretty",
	Paths:         []string{"./"},
	StopOnFailure: ut.Config.StopOnFailure,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	err := ut.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}
	if suite.Run() != 0 {
		log.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func AddRulesetCommand(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string, enabled string) error {
	dataProduct = ut.ProcessString(dataProduct)
	ruleset = ut.ProcessString(ruleset)
	commandString := "../gravity-cli product ruleset add "
	if dataProduct != testutils.NullString {
		commandString += " " + dataProduct
	}
	if ruleset != testutils.NullString {
		commandString += " " + ruleset
	}
	if event != testutils.IgnoreString {
		event := ut.ProcessString(event)
		commandString += " --event " + event
	}
	if method != testutils.IgnoreString {
		method := ut.ProcessString(method)
		commandString += " --method " + method
	}
	if pk != testutils.IgnoreString {
		pk := ut.ProcessString(pk)
		commandString += " --pk " + pk
	}
	if desc != testutils.IgnoreString {
		desc := ut.ProcessString(desc)
		commandString += " --desc " + desc
	}
	if handler != testutils.IgnoreString {
		commandString += " --handler " + handler
	}
	if schema != testutils.IgnoreString {
		commandString += " --schema " + schema
	}
	if enabled == testutils.TrueString {
		commandString += " --enabled"
	} else if enabled == testutils.FalseString {
		commandString += " --enabled=false"
	} else if enabled != testutils.IgnoreString {
		return errors.New("enabled parameter incorrect")
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func AddRulesetCommandFailed() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return fmt.Errorf("add ruleset should fail")
}

func AddRulesetCommandSuccess() error {
	if ut.CmdResult.Err == nil {
		return nil
	}
	return fmt.Errorf("add ruleset should success")
}

func SearchRulesetByCLISuccess(dataProduct string, ruleset string) error {
	dataProduct = ut.ProcessString(dataProduct)
	ruleset = ut.ProcessString(ruleset)
	cmd := exec.Command("../gravity-cli", "product", "ruleset", "info", dataProduct, ruleset, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	return err
}

func SearchRulesetByNatsSuccess(dataProduct string, ruleset string, method string, event string, pk string, desc string, handler string, schema string, enabled string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &jsonData)
	if err != nil {
		fmt.Println("error occurred while decoding JSON: ", err)
		return err
	}
	ruleset = ut.ProcessString(ruleset)
	method = ut.ProcessString(method)
	event = ut.ProcessString(event)
	desc = ut.ProcessString(desc)
	pk = ut.ProcessString(pk)

	if ruleset != testutils.NullString {
		if ruleset != jsonData.Rules[ruleset].Name {
			return errors.New("ruleset does not match the information in NATS")
		}
	}

	if err := ut.ValidateField(jsonData.Rules[ruleset].Method, method); err != nil {
		return err
	}
	if err := ut.ValidateField(jsonData.Rules[ruleset].Event, event); err != nil {
		return err
	}
	if err := ut.ValidateField(jsonData.Rules[ruleset].Desc, desc); err != nil {
		return err
	}
	pkStr := strings.Join(jsonData.Rules[ruleset].PrimaryKey, ",")
	if err := ut.ValidateField(pkStr, pk); err != nil {
		return err
	}

	if err := ut.ValidateHandler(jsonData.Rules[ruleset].Handler, handler); err != nil {
		return err
	}

	if err := ut.ValidateSchema(jsonData.Rules[ruleset].Schema, schema); err != nil {
		return err
	}

	if err := ut.ValidateEnabled(jsonData.Rules[ruleset].Enabled, enabled); err != nil {
		return err
	}
	return nil
}

// func AssertErrorMessages(expected string) error {
// 	if cmdResult.Stderr == expected {
// 		return nil
// 	}
// TODO
// 	return fmt.Errorf("the error message should be: %s", expected)
// }

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})
	ctx.Given(`^NATS has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product "'(.*?)'" and enabled is "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.When(`^"'(.*?)'" add ruleset "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, AddRulesetCommand)
	ctx.Then(`^CLI returns exit code 1$`, AddRulesetCommandFailed)
	ctx.Then(`^Check adding ruleset success$`, AddRulesetCommandSuccess)
	ctx.Then(`^Use gravity-cli to query the "'(.*?)'" "'(.*?)'" exists$`, SearchRulesetByCLISuccess)
	ctx.Then(`^Use NATS jetstream to query the "'(.*?)'" "'(.*?)'" exists and parameters "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" "'(.*?)'" are correct$`, SearchRulesetByNatsSuccess)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, AssertErrorMessages)
}
