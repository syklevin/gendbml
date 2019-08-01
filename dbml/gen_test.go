package dbml

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestGenModel(t *testing.T) {

	ba, err := ioutil.ReadFile("./KrcDB.dbml")
	if err != nil {
		t.Fatal(err)
	}

	dbml := &DBML{}

	err = xml.Unmarshal(ba, dbml)
	if err != nil {
		t.Fatal(err)
	}

	mg := NewModelGen(dbml, "data")

	err = mg.GenModelFile("models.go")
	if err != nil {
		t.Fatal(err)
	}

	err = mg.GenFuncFile("funcs.go")
	if err != nil {
		t.Fatal(err)
	}

}
