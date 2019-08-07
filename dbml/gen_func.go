package dbml

import (
	"fmt"
	"strings"
)

type FuncInfo struct {
	Name        string
	Inputs      []string
	Returns     []string
	FinalReturn string
	Body        string
}

func buildFuncInfo(fn DBMLFunc, pkg, externalDB string) (*FuncInfo, error) {
	if len(fn.DBMLFuncElements) > 1 {
		return nil, fmt.Errorf("not support gen multiple result set for %s", fn.Method)
	}

	fi := &FuncInfo{}
	fi.Name = strings.Title(fn.Method)
	fi.Inputs = []string{}
	fi.Inputs = append(fi.Inputs, "ctx context.Context")
	if len(externalDB) == 0 {
		fi.Inputs = append(fi.Inputs, "db *sqlx.DB")
	}

	if len(fn.Parameters) > 0 {
		fi.Inputs = append(fi.Inputs, "param *"+fi.Name+"Param")

		var fieldOut, fieldType string
		//append out params
		for _, p := range fn.Parameters {
			if p.Direction == "Out" || p.Direction == "InOut" {
				fieldOut = p.Name
				if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
					fieldOut = p.Name[1:]
				}
				if fieldOut == "RtnCode" || fieldOut == "RtnMessage" {
					continue
				}
				fieldOut = "out" + fieldOut
				fieldType = csTypeToGoType(p.Type, false)
				fieldOut = fieldOut + " *" + fieldType
				fi.Inputs = append(fi.Inputs, fieldOut)
			}
		}
	}

	var body strings.Builder
	fi.Returns = []string{}
	if len(fn.DBMLFuncElements) == 1 {
		el := fn.DBMLFuncElements[0]
		fi.Returns = append(fi.Returns, fmt.Sprintf("[]*%s", strings.Title(el.Name)))
		fi.FinalReturn = "&rst, "
		body.WriteString("var rst []" + strings.Title(el.Name) + "{}\n\t")
	}
	fi.Returns = append(fi.Returns, "err error")
	fi.FinalReturn += "checkError(err, errCode, errMsg)"

	var run string

	if len(externalDB) > 0 {
		if len(fn.DBMLFuncElements) == 0 {
			run = fmt.Sprintf("_, err = "+pkg+"."+externalDB+".ExecContext(ctx, \"%s\"", fn.Name)
		} else { //len(fn.DBMLFuncElements) == 1
			run = fmt.Sprintf("err = "+pkg+"."+externalDB+".SelectContext(ctx, &rst, \"%s\"", fn.Name)
		}
	} else {
		if len(fn.DBMLFuncElements) == 0 {
			run = fmt.Sprintf("_, err = db.ExecContext(ctx, \"%s\"", fn.Name)
		} else { //len(fn.DBMLFuncElements) == 1
			run = fmt.Sprintf("err = db.SelectContext(ctx, &rst, \"%s\"", fn.Name)
		}
	}

	body.WriteString(run)

	var paramFiledName string
	var namedfield string
	for _, p := range fn.Parameters {
		if p.Direction == "Out" || p.Direction == "InOut" {
			continue
		}
		paramFiledName = p.Name
		if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
			paramFiledName = p.Name[1:]
		}
		paramFiledName = "param." + strings.Title(paramFiledName)
		namedfield = fmt.Sprintf("\n\t\t,sql.Named(\"%s\", %s)", p.Name, paramFiledName)
		body.WriteString(namedfield)
	}

	for _, p := range fn.Parameters {
		if p.Direction == "Out" || p.Direction == "InOut" {
			paramFiledName = p.Name
			if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
				paramFiledName = p.Name[1:]
			}

			if paramFiledName == "RtnCode" {
				paramFiledName = "&errCode"
			} else if paramFiledName == "RtnMessage" {
				paramFiledName = "&errMsg"
			} else {
				paramFiledName = "out" + strings.Title(paramFiledName)
			}
			namedfield = fmt.Sprintf("\n\t\t,sql.Named(\"%s\", sql.Out{Dest: %s})", p.Name, paramFiledName)
			body.WriteString(namedfield)
		}
	}

	body.WriteString("\n\t)")

	fi.Body = body.String()

	return fi, nil

}
