package dbml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
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

	f, err := os.Create("models.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = mg.GenModelFile(f)
	if err != nil {
		t.Fatal(err)
	}

	f2, err := os.Create("funcs.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	err = mg.GenFuncFile(f2)
	if err != nil {
		t.Fatal(err)
	}

}
