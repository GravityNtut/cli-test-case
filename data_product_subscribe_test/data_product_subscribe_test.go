package dataproductsubscribe

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"test-case/testutils"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

var ut = testutils.TestUtils{}

func TestFeatures(t *testing.T) {
	ut.LoadConfig()
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

func SubscribeDataProduct() error {
	const TIMEOUT = 3
	cmd := exec.Command("../gravity-cli", "product", "sub", "drink", "-s", ut.Config.JetstreamURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()
	// 接收三秒內的資料
	after := time.After(TIMEOUT * time.Second)
	select {
	case <-after:
		cmd.Process.Signal(syscall.SIGINT)
		time.Sleep(10 * time.Millisecond)
		cmd.Process.Kill()
		break
	case <-done:
		break
	}
	fmt.Println(stdout.String())
	return nil
}

func DisplayData() error {
	return nil
}

func ValidateData() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.When(`^訂閱data product`, SubscribeDataProduct)
	ctx.Then(`^顯示資料`, DisplayData)
}
