package dbml

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type ModelGen struct {
	Pkg          string
	Cls          string
	DataPkg      string
	ErrorPkg     string
	ErrorPkgName string
	Models       []*ModelInfo
	Funcs        []*FuncInfo
	TestFuncs    []*TestFuncInfo
}

func NewModelGen(dbml *DBML, pkg, datapkg, externalDB, errorPkg string) *ModelGen {
	dataPkgList := strings.Split(datapkg, "/")
	dataPkgName := dataPkgList[len(dataPkgList)-1]

	errorPkgList := strings.Split(errorPkg, "/")
	errorpkgName := errorPkgList[len(errorPkgList)-1]

	mg := &ModelGen{
		Pkg:          pkg,
		Cls:          dbml.Class,
		DataPkg:      datapkg,
		ErrorPkg:     errorPkg,
		ErrorPkgName: errorpkgName,
		Models:       []*ModelInfo{},
		Funcs:        []*FuncInfo{},
	}

	for _, fn := range dbml.Functions {
		//input model
		param := buildParamModel(fn)
		mg.Models = append(mg.Models, param)
		results := buildResultModels(fn)
		mg.Models = append(mg.Models, results...)

		spfn, err := buildFuncInfo(fn, dataPkgName, externalDB)
		if err != nil {
			fmt.Println(err)
			continue
		}
		mg.Funcs = append(mg.Funcs, spfn)

		tsfn, err := buildTestFuncInfo(fn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		mg.TestFuncs = append(mg.TestFuncs, tsfn)

	}

	return mg
}

func (mg *ModelGen) GenModelFile(w io.Writer) error {

	t, err := template.New("Models template").Parse(ModelTmpl)
	if err != nil {
		return err
	}
	err = t.Execute(w, mg)
	return err
}

func (mg *ModelGen) GenFuncFile(w io.Writer) error {
	t, err := template.New("Funcs template").Parse(FuncTmpl)
	if err != nil {
		return err
	}
	err = t.Execute(w, mg)
	return err
}

func (mg *ModelGen) GenTestFuncFile(w io.Writer) error {
	t, err := template.New("Test Funcs template").Parse(TestFuncTmpl)
	if err != nil {
		return err
	}
	err = t.Execute(w, mg)
	return err
}

type ModelInfo struct {
	Name   string
	Fields []string
}

func buildParamModel(fn DBMLFunc) *ModelInfo {
	mi := &ModelInfo{}
	mi.Name = strings.Title(fn.Method) + "Param"
	mi.Fields = genParamModelFields(fn.Parameters)
	return mi
}

func buildResultModels(fn DBMLFunc) []*ModelInfo {
	ret := []*ModelInfo{}
	for _, el := range fn.DBMLFuncElements {
		mi := &ModelInfo{}
		mi.Name = strings.Title(el.Name)
		mi.Fields = genResultModelFields(el.Columns)
		ret = append(ret, mi)
	}
	return ret
}

func genParamModelFields(params []DBMLFuncParam) []string {
	arr := []string{}
	var field string
	var anno string
	for _, p := range params {
		if p.Direction == "Out" || p.Direction == "InOut" {
			continue
		}
		if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
			field = p.Name[1:]
		} else {
			field = p.Name
		}
		field = strings.Title(field)
		anno = fmt.Sprintf("db:\"%s\"", p.Name)
		arr = append(arr, fmt.Sprintf("%s %s `%s`", field, csTypeToGoType(p.Type, false), anno))
	}
	return arr
}

func genResultModelFields(cols []DBMLFuncElementColumn) []string {
	arr := []string{}
	var field string
	var anno string
	var nullable bool
	for _, p := range cols {
		nullable = true
		if p.CanBeNull == "false" {
			nullable = false
		}
		if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
			field = p.Name[1:]
		} else {
			field = p.Name
		}
		field = strings.Title(field)
		anno = fmt.Sprintf("db:\"%s\"", p.Name)
		arr = append(arr, fmt.Sprintf("%s %s `%s`", field, csTypeToGoType(p.Type, nullable), anno))
	}
	return arr
}

func csTypeToGoType(csType string, nullable bool) string {
	switch csType {
	case "System.Byte", "System.Int32", "System.Int16":
		if nullable {
			return "sql.NullInt64"
		}
		return "int"
	case "System.Int64":
		if nullable {
			return "sql.NullInt64"
		}
		return "int64"
	case "System.Char", "System.String", "System.Xml.Linq.XElement":
		if nullable {
			return "sql.NullString"
		}
		return "string"
	case "System.Boolean":
		if nullable {
			return "sql.NullBool"
		}
		return "bool"
	case "System.DateTime":
		if nullable {
			return "*time.Time"
		}
		return "time.Time"
	case "System.Double", "System.Decimal":
		if nullable {
			return "sql.NullFloat64"
		}
		return "float64"
	}

	return ""
}
