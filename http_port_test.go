package main

import (
	"os"
	"testing"
)

func TestHttpDisplayURL(t *testing.T) {
	original := HttpDisplayPort
	t.Cleanup(func() { HttpDisplayPort = original })

	HttpDisplayPort = 10100
	if got := httpDisplayURL("/web"); got != "http://127.0.0.1:10100/web" {
		t.Fatalf("httpDisplayURL(/web) = %q", got)
	}
	if got := httpDisplayURL("dm"); got != "http://127.0.0.1:10100/dm" {
		t.Fatalf("httpDisplayURL(dm) = %q", got)
	}
	if got := musicServerURL("/music"); got != "http://127.0.0.1:99/music" {
		t.Fatalf("musicServerURL = %q", got)
	}
}

func TestHttpDisplayListenPortsDefault(t *testing.T) {
	t.Setenv("BILILINE_HTTP_PORT", "")
	got := httpDisplayListenPorts()
	if len(got) != 2 || got[0] != preferredHttpDisplayPort || got[1] != fallbackHttpDisplayPort {
		t.Fatalf("default ports = %#v", got)
	}
}

func TestHttpDisplayListenPortsOverride(t *testing.T) {
	t.Setenv("BILILINE_HTTP_PORT", "18080")
	got := httpDisplayListenPorts()
	if len(got) != 1 || got[0] != 18080 {
		t.Fatalf("override ports = %#v", got)
	}
}

func TestHttpDisplayListenPortsInvalidOverride(t *testing.T) {
	t.Setenv("BILILINE_HTTP_PORT", "not-a-port")
	got := httpDisplayListenPorts()
	if len(got) != 2 || got[0] != preferredHttpDisplayPort {
		t.Fatalf("invalid override ports = %#v", got)
	}
	_ = os.Unsetenv("BILILINE_HTTP_PORT")
}
