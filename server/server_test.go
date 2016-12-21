package server

import "testing"

var s Server

func TestNew(t *testing.T) {
	s = New()
}

func TestRun(t *testing.T) {
	s.Run()
}
