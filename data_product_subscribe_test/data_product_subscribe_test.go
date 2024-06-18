package dataproductsubscribe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"test-case/testutils"
	"testing"

	"github.com/cucumber/godog"
)

var ut = testutils.TestUtils{}
var EventCount = 10

type JSONSubData struct {
	Event     string      `json:"event"`
	Header    interface{} `json:"header"`
	Method    string      `json:"method"`
	Payload   interface{} `json:"payload"`
	Pk        []string    `json:"primarykey"`
	Product   string      `json:"product"`
	Seq       uint64      `json:"seq"`
	Subject   string      `json:"subject"`
	Table     string      `json:"table"`
	Timestamp string      `json:"timestamp"`
}

func TestFeatures(t *testing.T) {
	err := ut.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"./"},
			StopOnFailure: ut.Config.StopOnFailure,
			TestingT:      t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
func FindJSON(data string) []string {
	var result []string
	stringStart := -1
	lvl := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '{' {
			if lvl == 0 {
				stringStart = i
			}
			lvl++
		} else if data[i] == '}' {
			lvl--
			if lvl == 0 {
				resultString := data[stringStart : i+1]
				result = append(result, resultString)
			}
		}
	}
	return result
}

var Timeout int

func SetTimeout(timeout int) error {
	Timeout = timeout
	return nil
}

func SubscribeDataProductCommand(productName string, _ string, partitions string, seq string) error {
	productName = ut.ProcessString(productName)
	cmdString := fmt.Sprintf("timeout %d ../gravity-cli product sub ", Timeout)
	if productName != testutils.NullString {
		cmdString += productName + " "
	}
	if partitions != testutils.IgnoreString {
		cmdString += "--partitions " + partitions + " "
	}
	if seq != testutils.IgnoreString {
		cmdString += "--seq " + seq + " "
	}

	cmdString += "-s " + ut.Config.JetstreamURL
	err := ut.ExecuteShell(cmdString)
	if err != nil {
		return err
	}
	return nil
}

func DisplayData() error {
	resultStringList := FindJSON(ut.CmdResult.Stdout)
	fmt.Println(resultStringList)
	return nil
}

func PublishProductEvent() error {
	EventCount = 10
	// Generate 7 events which id is 87
	for i := 1; i <= EventCount-3; i++ {
		result := fmt.Sprintf(`{"id":87, "name":"test%d", "kcal":%d, "price":%d}`, i, i*100, i+20)
		cmd := exec.Command(testutils.GravityCliString, "pub", Event, result, "-s", ut.Config.JetstreamURL)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	// Generate 3 events which id is 99
	for i := EventCount - 2; i <= EventCount; i++ {
		result := fmt.Sprintf(`{"id":99, "name":"test%d", "kcal":%d, "price":%d}`, i, i*300, i+30)
		cmd := exec.Command(testutils.GravityCliString, "pub", Event, result, "-s", ut.Config.JetstreamURL)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("publish event failed: %s", err.Error())
		}
	}
	return nil
}

func CreateDataProduct(dataProduct string) error {
	decription := "drink data"
	schema := "./assets/schema.json"
	enabled := "true"
	cmd := exec.Command(testutils.GravityCliString, "product", "create", dataProduct, "--desc", decription, "--schema", schema, "--enabled="+enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("create data product failed")
	}
	return nil
}

var DataProduct string
var Ruleset string
var Event string
var ExpectPartition []string

const Method string = "create"
const Pk string = "id"
const RulesetDesc string = "drink create event"
const Handler string = "./assets/handler.js"
const Schema string = "./assets/schema.json"
const Enabled string = "true"

func CreateDataProductRuleset(dataProduct string, ruleset string) error {
	DataProduct = dataProduct
	Ruleset = ruleset
	Event = ruleset
	cmd := exec.Command(testutils.GravityCliString, "product", "ruleset", "add", DataProduct, Ruleset, "--event", Event, "--method", Method, "--desc", RulesetDesc, "--pk", Pk, "--handler", Handler, "--schema", Schema, "--enabled="+Enabled, "-s", ut.Config.JetstreamURL)
	err := cmd.Run()
	if err != nil {
		return errors.New("data product add ruleset failed")
	}
	return nil
}

func ValidateSubResult(partitionString string, seqString string) error {
	var err error
	var seq uint64 = 1
	if seqString != testutils.IgnoreString {
		seq, err = strconv.ParseUint(seqString, 10, 64)
		if err != nil {
			return fmt.Errorf("fail from covert seq to uint: %s", err.Error())
		}
	}
	resultStringList := FindJSON(ut.CmdResult.Stdout)
	numbersOfEvents := uint64(EventCount) - seq + 1
	if uint64(EventCount) < seq {
		numbersOfEvents = 0
	}

	if partitionString == "131" {
		if numbersOfEvents > 3 {
			numbersOfEvents -= 3
		} else {
			numbersOfEvents = 0
		}
	} else if partitionString == "200" {
		if numbersOfEvents > 3 {
			numbersOfEvents = 3
		}
	} else if partitionString != "131,200" && partitionString != "-1" && partitionString != testutils.IgnoreString {
		numbersOfEvents = 0
	}

	if uint64(len(resultStringList)) != numbersOfEvents {
		return fmt.Errorf("the number of retrieved events differs from the number of published events. Expect: %d, Retrieve: %d", numbersOfEvents, len(resultStringList))
	}
	for i, jsonData := range resultStringList {
		i := uint64(i)

		if partitionString == "200" {
			if seq < 8 {
				i = 8 + i
			} else {
				i = i + seq
			}
		} else {
			i = i + seq
		}

		var UnmarshalResult JSONSubData
		if err := json.Unmarshal([]byte(jsonData), &UnmarshalResult); err != nil {
			return errors.New("json unmarshal failed" + err.Error())
		}

		var payloadString string
		if i >= uint64(EventCount)-3 {
			payloadString = fmt.Sprintf(`{"id":99, "name":"test%d", "kcal":%d, "price":%d}`, i, i*300, i+30)
		} else {
			payloadString = fmt.Sprintf(`{"id":87, "name":"test%d", "kcal":%d, "price":%d}`, i, i*100, i+20)
		}
		expectPayload := ut.FormatJSONData(payloadString)

		if err := ut.AssertStringEqual(UnmarshalResult.Event, Event); err != nil {
			return err
		}
		if err := ut.AssertStringEqual(UnmarshalResult.Product, DataProduct); err != nil {
			return err
		}
		payloadByte, _ := json.Marshal(UnmarshalResult.Payload)
		resultPayload := string(payloadByte)
		if err := ut.AssertStringEqual(resultPayload, expectPayload); err != nil {
			return err
		}
		pkExpect := strings.Split(Pk, ",")
		for i := 0; i < len(pkExpect); i++ {
			if err := ut.AssertStringEqual(UnmarshalResult.Pk[i], pkExpect[i]); err != nil {
				return err
			}
		}
		// The seq is to be modified in the future, so we will not check it for now.
		// if err := ut.AssertUIntEqual(UnmarshalResult.Seq, i+1); err != nil {
		// 	return err
		// }
		if err := ut.AssertStringEqual(UnmarshalResult.Method, Method); err != nil {
			return err
		}
	}
	return nil
}

func SubCommandFail() error {
	if ut.CmdResult.Err == nil || ut.CmdResult.Err.(*exec.ExitError).ExitCode() == 124 {
		return errors.New("this command should be failed")
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ut.ClearDataProducts()
		return ctx, nil
	})

	ctx.Given(`^NATS has been opened$`, ut.CheckNatsService)
	ctx.Given(`^Dispatcher has been opened$`, ut.CheckDispatcherService)
	ctx.Given(`^Create data product "'([^'"]*?)'"$`, CreateDataProduct)
	ctx.Given(`^Create data product "'([^'"]*?)'" with ruleset "'([^'"]*?)'"$`, CreateDataProductRuleset)
	ctx.Given(`^Publish 10 Events$`, PublishProductEvent)
	ctx.Given(`^Set subscribe timeout with "'(.*?)'" second$`, SetTimeout)
	ctx.When(`^Subscribe data product "'(.*?)'" using parameters "'(.*?)'" "'(.*?)'" "'(.*?)'"`, SubscribeDataProductCommand)
	ctx.Then(`^The CLI returns all events data within the "'(.*?)'" and after "'(.*?)'"$`, ValidateSubResult)
	ctx.Then(`^The CLI returns exit code 1$`, SubCommandFail)
}
