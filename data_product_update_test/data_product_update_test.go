package dataproductupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

type JSONData struct {
	Name    string      `json:"name"`
	Desc    string      `json:"desc"`
	Enabled bool        `json:"enabled"`
	Schema  interface{} `json:"schema"`
}

var newJSONData JSONData
var ut testutils.TestUtils
var originJSONData string

var opts = godog.Options{
	Format:        "pretty",
	Paths:         []string{"./"},
	StopOnFailure: ut.Config.StopOnFailure,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(_ *testing.M) {
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

func UpdateDataProductCommand(dataProduct string, description string, enabled string, schema string) error {
	commandString := "../gravity-cli product update "
	if dataProduct != testutils.NullString {
		commandString += dataProduct
	}
	if description != testutils.IgnoreString {
		commandString += " --desc " + description
	}

	if enabled == testutils.TrueString {
		commandString += " --enabled"
	} else if enabled == testutils.FalseString {
		commandString += " --enabled=false"
	} else if enabled != testutils.IgnoreString {
		return errors.New("enabled parameter incorrect")
	}

	if schema != testutils.IgnoreString {
		commandString += " --schema " + schema
	}
	commandString += " -s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(commandString)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDataProductCommandFail() error {
	if ut.CmdResult.Err != nil {
		return nil
	}
	return errors.New("update data product fail")
}

func SearchDataProductByJetstreamSuccess(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
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
	return errors.New("the data product was not created successfully in jetstream")
}

// func AssertErrorMessages(errorMessage string) error {
// 	outErr := cmdResult.stderr
// 	if outErr == errorMessage {
// 		return nil
// 	}
// 	return errors.New("the error message should be: %s", errorMessage)
// }

func DataProductNotChanges(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)

	if string(entry.Value()) == originJSONData {
		return nil
	}
	return errors.New("does not match the original data")
}

func DataProductUpdateSuccess(dataProduct string, desc string, schema string, enabled string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	err = json.Unmarshal((entry.Value()), &newJSONData)
	if err != nil {
		fmt.Println("error occurred while decoding JSON: ", err)
		return err
	}

	if err := ut.ValidateSchema(newJSONData.Schema, schema); err != nil {
		return err
	}

	if err := ut.ValidateEnabled(newJSONData.Enabled, enabled); err != nil {
		return err
	}

	if err := ut.ValidateField(newJSONData.Desc, desc); err != nil {
		return err
	}

	if dataProduct != newJSONData.Name {
		return errors.New("update data error")
	}
	return nil
}

func UpdateDataProductCommandSuccess() error {
	if ut.CmdResult.Err != nil {
		return errors.New("update data fail")
	}
	return nil
}

func StoreNowDataProduct(dataProduct string) error {
	nc, _ := nats.Connect(testutils.NatsProtocol + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	kv, _ := js.KeyValue("GVT_default_PRODUCT")
	entry, _ := kv.Get(dataProduct)
	originJSONData = string(entry.Value())
	return nil
}

// TODO
// func AssertErrorMessages(errorMessage string) error {

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^NATS has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product "'(.*?)'" and enabled is "'(.*?)'"$`, ut.CreateDataProduct)
	ctx.Given(`^Store NATS copy of existing data product "'(.*?)'"$`, StoreNowDataProduct)
	ctx.When(`^Update the name of data product to "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, UpdateDataProductCommand)
	ctx.Then(`^Check updating data product success$`, UpdateDataProductCommandSuccess)
	ctx.Then(`^CLI returns exit code 1$`, UpdateDataProductCommandFail)
	// ctx.Then(`^The error message should be "'(.*?)'"$`, AssertErrorMessages)
	ctx.Then(`^Use NATS jetstream to query the "'(.*?)'" update successfully and parameters are "'(.*?)'" "'(.*?)'" "'(.*?)'"$`, DataProductUpdateSuccess)
	ctx.Then(`^Use NATS jetstream to query the "'(.*?)'" without changing parameters$`, DataProductNotChanges)
}
