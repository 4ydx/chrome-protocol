package cdp

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	f, err := os.Create("testing.log")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(f)

	res := m.Run()
	os.Exit(res)
}
