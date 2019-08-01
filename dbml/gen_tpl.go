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
    _ = sql.Named{}
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
