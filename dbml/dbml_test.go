package dbml

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestDbmlUnmarshal(t *testing.T) {

	ba, err := ioutil.ReadFile("./KrcDB.dbml")
	if err != nil {
		t.Fatal(err)
	}

	var dbml DBML

	err = xml.Unmarshal(ba, &dbml)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(dbml)
}
