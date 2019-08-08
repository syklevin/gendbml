package dbml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenModel(t *testing.T) {

	ba, err := ioutil.ReadFile("../fixtures/test.dbml")
	if err != nil {
		t.Fatal(err)
	}

	dbml := &DBML{}

	err = xml.Unmarshal(ba, dbml)
	if err != nil {
		t.Fatal(err)
	}

	mg := NewModelGen(dbml, "data", "github.acsdev.net/helios/zeus/pkg/data", "GameClub", "github.acsdev.net/helios/zeus/pkg/errors")

	outDir := "../tmp"

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.Mkdir(outDir, 0755)
	}

	fp := filepath.Join(outDir, "models.go")

	f, err := os.Create(fp)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = mg.GenModelFile(f)
	if err != nil {
		t.Fatal(err)
	}

	fp = filepath.Join(outDir, "funcs.go")

	f2, err := os.Create(fp)
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	err = mg.GenFuncFile(f2)
	if err != nil {
		t.Fatal(err)
	}

	fp = filepath.Join(outDir, "funcs_test.go")

	f3, err := os.Create(fp)
	if err != nil {
		t.Fatal(err)
	}
	defer f3.Close()

	err = mg.GenTestFuncFile(f3)
	if err != nil {
		t.Fatal(err)
	}

}
