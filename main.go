package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/syklevin/gendbml/dbml"
)

func die(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

func main() {

	var dbmlFile string
	var outDir string
	var pkgName string
	flag.StringVar(&dbmlFile, "file", "", "dbml file for gen")
	flag.StringVar(&pkgName, "pkg", "data", "output pkg namen")
	flag.StringVar(&outDir, "dir", "", "output dir for gen")
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

	mg := dbml.NewModelGen(ml, pkgName)

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
}
