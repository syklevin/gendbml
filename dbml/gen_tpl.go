package dbml

var ModelTmpl = `package {{.Pkg}}

// gen from dbml {{.Cls}}


import (
    "database/sql"
    "time"
)
var (
    _ = time.Second
    _ = sql.LevelDefault
)

{{range .Models}}
type {{.Name}} struct {
    {{range .Fields}}{{.}}
    {{end}}
}
{{end}}
`

var FuncTmpl = `package {{.Pkg}}

// gen from dbml {{.Cls}}

import (
    "context"
    "database/sql"
    "github.com/jmoiron/sqlx"
    
    {{if .DataPkg}}"{{.DataPkg}}" {{end}}
    {{if .ErrorPkg}}"{{.ErrorPkg}}" {{end}}
)


var (
    _ = context.Background()
    _ = sql.LevelDefault
    _ = sqlx.DB{}
)

{{if .ErrorPkgName}}func checkError(err error, errCode int32, errMsg string) error {
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if errCode != 0 {
		return {{.ErrorPkgName}}.Coded(int(errCode), errMsg)
	}
	return err
}
{{end}}

{{range .Funcs}}
func {{.Name}}(
    {{- range $i, $e := .Inputs -}}
        {{if $i}}, {{end}}{{$e}}
    {{- end}}
    {{- $length := len .Returns}}
        {{- if gt $length 1}}) ({{else}}) {{end}}
    {{- range $j, $k := .Returns -}}
        {{if $j}}, {{end}}{{$k}}
    {{- end}}
    {{- if gt $length 1}}) { {{- else}} { {{- end}}
    var errCode int32
    var errMsg string

    {{.Body}}
    return {{.FinalReturn}}
}
{{end}}

`

var TestFuncTmpl = `package {{.Pkg}}

import (
    "testing"
    "context"
    "github.com/jmoiron/sqlx"
)

var (
    db_ *sqlx.DB = nil //change this!!
)

{{range .TestFuncs}}
func Test{{.Name}}(t *testing.T){
    {{range .Declares}}{{.}}
    {{end}}
    {{- range $j, $k := .Returns -}}
        {{if $j}}, {{end}}{{$k}}
    {{- end}} = {{.Name}}(
    {{- range $i, $e := .Arguments -}}
        {{if $i}}, {{end}}{{$e}}
    {{- end}})
    if err != nil {
        t.Fatal(err)
    }
    {{ if eq (len .Returns) 2 }}
    if rst == nil {
        t.Fatal("rst should not be nil")
    }
    {{ end }}
}
{{end}}


`
