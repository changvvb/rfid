package server

import "testing"

func TestNew(t *testing.T) {
	server := New(6060)
	t.Log("test port: 6060")
	if server.port != 6060 {
		t.Error("port error")
	}
}

func TestRun(t *testing.T) {
	server := New(6060)
	server.Run()
}
