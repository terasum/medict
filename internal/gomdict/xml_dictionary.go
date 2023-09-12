package gomdict

import (
	"encoding/xml"
	"strings"
)

// Dictionary was generated 2023-09-11 11:07:50 by https://xml-to-go.github.io/ in Ukraine.
type Dictionary struct {
	XMLName                  xml.Name `xml:"Dictionary"`
	Text                     string   `xml:",chardata"`
	GeneratedByEngineVersion string   `xml:"GeneratedByEngineVersion,attr"`
	RequiredEngineVersion    string   `xml:"RequiredEngineVersion,attr"`
	Encrypted                string   `xml:"Encrypted,attr"`
	Encoding                 string   `xml:"Encoding,attr"`
	Format                   string   `xml:"Format,attr"`
	Stripkey                 string   `xml:"Stripkey,attr"`
	CreationDate             string   `xml:"CreationDate,attr"`
	Compact                  string   `xml:"Compact,attr"`
	Compat                   string   `xml:"Compat,attr"`
	KeyCaseSensitive         string   `xml:"KeyCaseSensitive,attr"`
	Description              string   `xml:"Description,attr"`
	Title                    string   `xml:"Title,attr"`
	DataSourceFormat         string   `xml:"DataSourceFormat,attr"`
	StyleSheet               string   `xml:"StyleSheet,attr"`
	Left2Right               string   `xml:"Left2Right,attr"`
	RegisterBy               string   `xml:"RegisterBy,attr"`
}

func parseXMLHeader(xmldata string) (*Dictionary, error) {
	dic := &Dictionary{}
	err := xml.Unmarshal([]byte(strings.TrimSpace(xmldata)), dic)
	if err != nil {
		return nil, err
	}
	return dic, nil
}
