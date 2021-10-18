package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	// disable log output for testing
	log.SetOutput(ioutil.Discard)
	m.Run()
}
