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
)


var (
    _ = context.Background()
    _ = sql.LevelDefault
    _ = sqlx.DB{}
)

{{range .Funcs}}
func {{.Name}}(
    {{- range $i, $e := .Inputs -}}
        {{if $i}}, {{end}}{{$e}}
    {{- end}}) (
    {{- range $j, $k := .Returns -}}
        {{if $j}}, {{end}}{{$k}}
    {{- end}}) {
    {{.Body}}
    return
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
