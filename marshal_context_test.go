package gosnmp_test

import (
	"context"
	"errors"
	g "github.com/gosnmp/gosnmp"
	"net"
	"testing"
	"time"
)

func TestCancelReceiving(t *testing.T) {

	srvr, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		t.Fatalf("udp4 error listening: %s", err)
	}
	defer srvr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	g.Default.Target = srvr.LocalAddr().(*net.UDPAddr).IP.String()
	g.Default.Port = uint16(srvr.LocalAddr().(*net.UDPAddr).Port)
	g.Default.Timeout = 25 * time.Second
	g.Default.Retries = 0
	g.Default.Context = ctx

	err = g.Default.Connect()
	if err != nil {
		t.Fatalf("Connect(%s) err: %v", g.Default.Target, err)
	}
	defer g.Default.Conn.Close()

	_, err = g.Default.Get([]string{"1.3.6.1.2.1.1.4.0"})
	if err == nil {
		t.Fatalf("Get() err should return error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %q", err.Error())
	}

}
