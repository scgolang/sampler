package sampler_test

import (
	"net"
	"testing"

	"github.com/scgolang/osc"
	"github.com/scgolang/sampler"
)

func TestAdd(t *testing.T) {
	const scsynthAddr = "127.0.0.1:57120"

	_ = mockScsynth(t, scsynthAddr)
	samps := newTestSampler(t, scsynthAddr)
	if err := samps.Add("MD16_Cow_01.wav", 0); err != nil {
		t.Fatal(err)
	}
}

func newTestSampler(t *testing.T, scsynthAddr string) *sampler.Sampler {
	s, err := sampler.New(scsynthAddr)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func mockScsynth(t *testing.T, listenAddr string) osc.Conn {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := osc.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	synthdefDoneMsg := osc.Message{
		Address: "/done",
		Arguments: osc.Arguments{
			osc.String("/d_recv"),
		},
	}
	go func() {
		if err := conn.Serve(1, osc.Dispatcher{
			"/d_recv": osc.Method(func(m osc.Message) error {
				return conn.SendTo(m.Sender, synthdefDoneMsg)
			}),
			"/g_new": osc.Method(func(m osc.Message) error {
				return nil
			}),
		}); err != nil {
			t.Fatal(err)
		}
	}()
	return conn
}
