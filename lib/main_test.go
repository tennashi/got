package got_test

import (
	"os"
	"path/filepath"
	"testing"

	got "github.com/tennashi/got/lib"
)

var testConfig map[string]*got.Config

func TestMain(m *testing.M) {
	beforeAll()
	code := m.Run()
	os.Exit(code)
}

func beforeAll() {
	initTestConfig()
}

func initTestConfig() {
	testConfig = map[string]*got.Config{
		"empty": {},
		"valid": {},
	}
	got.CaddPath(testConfig["valid"], filepath.Join("..", "testdata", "config_valid.golden"))
	got.Cload(testConfig["valid"])
}

func resetTestConfig() {
	initTestConfig()
}
