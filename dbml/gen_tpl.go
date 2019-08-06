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

func checkError(err error, errCode int32, errMsg string) error {
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if errCode != 0 {
		return {{.ErrorPkgName}}.Coded(int(errCode), errMsg)
	}
	return err
}

{{range .Funcs}}
func {{.Name}}(
    {{- range $i, $e := .Inputs -}}
        {{if $i}}, {{end}}{{$e}}
    {{- end}}) (
    {{- range $j, $k := .Returns -}}
        {{if $j}}, {{end}}{{$k}}
    {{- end}}) {
    var errCode int32
    var errMsg string

    {{.Body}}
    return {{.FinalReturn}}
}
{{end}}

`
