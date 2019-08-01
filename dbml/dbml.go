package dbml

import (
	"encoding/xml"
)

type DBML struct {
	XMLName   xml.Name   `xml:"Database"`
	Name      string     `xml:"Name,attr"`
	Class     string     `xml:"Class,attr"`
	Functions []DBMLFunc `xml:"Function"`
}

type DBMLFunc struct {
	Name             string            `xml:"Name,attr"`
	Method           string            `xml:"Method,attr"`
	Parameters       []DBMLFuncParam   `xml:"Parameter"`
	DBMLFuncElements []DBMLFuncElement `xml:"ElementType"`
	Return           DBMLFuncReturn    `xml:"Return"`
}

type DBMLFuncParam struct {
	Name      string `xml:"Name,attr"`
	Type      string `xml:"Type,attr"`
	DbType    string `xml:"DbType,attr"`
	Direction string `xml:"Direction,attr,omitempty"`
}

type DBMLFuncReturn struct {
	Type string `xml:"Type,attr"`
}

type DBMLFuncElement struct {
	Name    string                  `xml:"Name,attr"`
	Columns []DBMLFuncElementColumn `xml:"Column"`
}

type DBMLFuncElementColumn struct {
	Name      string `xml:"Name,attr"`
	Type      string `xml:"Type,attr"`
	DbType    string `xml:"DbType,attr"`
	CanBeNull string `xml:"CanBeNull,attr"`
}
