package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/NatalieACS/gendbml/dbml"
	//"github.com/syklevin/gendbml/dbml"
)

func die(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

func main() {

	var dbmlFile string
	var outDir string
	var pkgName string
	var dataPkg string
	var externalDB string
	var errorPkg string
	flag.StringVar(&dbmlFile, "file", "", "dbml file for gen")
	flag.StringVar(&pkgName, "pkg", "data", "output pkg namen")
	flag.StringVar(&outDir, "dir", "", "output dir for gen")
	flag.StringVar(&dataPkg, "datapkg", "", "pkg where singalten DB(s) appear")
	flag.StringVar(&externalDB, "externalDB", "", "DB variable name used to call DB in data pkg")
	flag.StringVar(&errorPkg, "errorPkg", "", "pkg where custumized error to be used")
	flag.Parse()

	ba, err := ioutil.ReadFile(dbmlFile)
	if err != nil {
		die(err)
	}

	ml := &dbml.DBML{}

	err = xml.Unmarshal(ba, ml)
	if err != nil {
		die(err)
	}

	mg := dbml.NewModelGen(ml, pkgName, dataPkg, externalDB, errorPkg)

	if outDir == "" {
		outDir = pkgName
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.Mkdir(outDir, 0755)
	}

	clsName := mg.Cls
	if idx := strings.Index(mg.Cls, "DataContext"); idx > 0 {
		clsName = mg.Cls[:idx]
	}

	modelFile := filepath.Join(outDir, strings.ToLower(clsName)+"_models.go")
	f, err := os.Create(modelFile)
	if err != nil {
		die(err)
	}
	defer f.Close()

	err = mg.GenModelFile(f)
	if err != nil {
		die(err)
	}

	funcFile := filepath.Join(outDir, strings.ToLower(clsName)+"_funcs.go")
	f2, err := os.Create(funcFile)
	if err != nil {
		die(err)
	}
	defer f2.Close()

	err = mg.GenFuncFile(f2)
	if err != nil {
		die(err)
	}

	testfuncFile := filepath.Join(outDir, strings.ToLower(clsName)+"_funcs_test.go")
	f3, err := os.Create(testfuncFile)
	if err != nil {
		die(err)
	}
	defer f3.Close()

	err = mg.GenTestFuncFile(f3)
	if err != nil {
		die(err)
	}

}
