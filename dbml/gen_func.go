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

func buildFuncInfo(fn DBMLFunc, dataPkgName, externalDB string) (*FuncInfo, error) {
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
		fi.FinalReturn = "rst, "
		body.WriteString("rst := []" + strings.Title(el.Name) + "{}\n\t")
	}
	fi.Returns = append(fi.Returns, "error")
	fi.FinalReturn += "checkError(err, errCode, errMsg)"

	var run string

	if len(externalDB) > 0 {
		if len(fn.DBMLFuncElements) == 0 {
			run = fmt.Sprintf("_, err = "+dataPkgName+"."+externalDB+".ExecContext(ctx, \"%s\",", fn.Name)
		} else { //len(fn.DBMLFuncElements) == 1
			run = fmt.Sprintf("err = "+dataPkgName+"."+externalDB+".SelectContext(ctx, &rst, \"%s\",", fn.Name)
		}
	} else {
		if len(fn.DBMLFuncElements) == 0 {
			run = fmt.Sprintf("_, err = db.ExecContext(ctx, \"%s\",", fn.Name)
		} else { //len(fn.DBMLFuncElements) == 1
			run = fmt.Sprintf("err = db.SelectContext(ctx, &rst, \"%s\",", fn.Name)
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
		namedfield = fmt.Sprintf("\n\t\tsql.Named(\"%s\", %s),", p.Name, paramFiledName)
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
			namedfield = fmt.Sprintf("\n\t\tsql.Named(\"%s\", sql.Out{Dest: %s}),", p.Name, paramFiledName)
			body.WriteString(namedfield)
		}
	}

	body.WriteString("\n\t)")

	fi.Body = body.String()

	return fi, nil

}

type TestFuncInfo struct {
	Name      string
	Declares  []string
	Arguments []string
	Returns   []string
	Body      string
}

func buildTestFuncInfo(fn DBMLFunc) (*TestFuncInfo, error) {

	if len(fn.DBMLFuncElements) > 1 {
		return nil, fmt.Errorf("not support gen multiple result set for %s", fn.Method)
	}

	fi := &TestFuncInfo{}
	fi.Name = strings.Title(fn.Method)
	fi.Declares = []string{}
	fi.Arguments = []string{}

	fi.Declares = append(fi.Declares, "var db = db_")
	fi.Declares = append(fi.Declares, "var ctx = context.Background()")
	fi.Arguments = append(fi.Arguments, "ctx", "db")

	if len(fn.Parameters) > 0 {
		fi.Declares = append(fi.Declares, "var param = &"+fi.Name+"Param{}")
		fi.Arguments = append(fi.Arguments, "param")

		var fieldOut, fieldType string
		//append out params
		for _, p := range fn.Parameters {
			if p.Direction == "Out" || p.Direction == "InOut" {
				fieldOut = p.Name
				if strings.Index(p.Name, "p") == 0 {
					fieldOut = p.Name[1:]
				}
				fieldOut = "out" + fieldOut
				fieldType = csTypeToGoType(p.Type, false)
				fi.Declares = append(fi.Declares, "var "+fieldOut+" "+fieldType)
				fi.Arguments = append(fi.Arguments, "&"+fieldOut)
			}
		}
	}

	fi.Returns = []string{}
	if len(fn.DBMLFuncElements) == 1 {
		el := fn.DBMLFuncElements[0]
		fi.Declares = append(fi.Declares, fmt.Sprintf("var rst = []*%s{}", strings.Title(el.Name)))
		fi.Returns = append(fi.Returns, "rst")
	}
	fi.Declares = append(fi.Declares, "var err error")
	fi.Returns = append(fi.Returns, "err")

	var body strings.Builder
	var run string
	if len(fn.DBMLFuncElements) == 0 {
		run = fmt.Sprintf("_, err = db.ExecContext(ctx, \"%s\"", fn.Name)
	} else { //len(fn.DBMLFuncElements) == 1
		run = fmt.Sprintf("err = db.SelectContext(ctx, &rst, \"%s\"", fn.Name)
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
		namedfield = fmt.Sprintf("\n\t\tsql.Named(\"%s\", %s),", p.Name, paramFiledName)
		body.WriteString(namedfield)
	}

	for _, p := range fn.Parameters {
		if p.Direction == "Out" || p.Direction == "InOut" {
			paramFiledName = p.Name
			if strings.Index(p.Name, "p") == 0 || strings.Index(p.Name, "w") == 0 {
				paramFiledName = p.Name[1:]
			}
			paramFiledName = "out" + strings.Title(paramFiledName)
			namedfield = fmt.Sprintf("\n\t\tsql.Named(\"%s\", sql.Out{Dest: %s}),", p.Name, paramFiledName)
			body.WriteString(namedfield)
		}
	}

	body.WriteString(")")

	fi.Body = body.String()

	return fi, nil

}
