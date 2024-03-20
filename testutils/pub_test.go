package testutils

import (
	"os"
	"testing"

	"github.com/nats-io/nats.go"
)

var ut TestUtils

func TestMyFunction(t *testing.T) {
	ut.LoadConfig()
	nc, _ := nats.Connect("nats://" + ut.Config.JetstreamURL)
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Error(err)
	}
	t.Log("Connected to JetStream")

	// js.Subscribe("$GVT.default.EVENT.*", func(m *nats.Msg) {
	// 	t.Log("Received message", string(m.Data))
	// }, nats.AckAll())

	f, _ := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	js.Subscribe("$GVT.default.DP.abc.0.EVENT.*", func(m *nats.Msg) {
		t.Log("Received message", string(m.Data))
	}, nats.AckAll())
}
