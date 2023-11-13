//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gomdict

import (
	"encoding/xml"
	"strings"
)

// Dictionary was generated 2023-09-11 11:07:50 by https://xml-to-go.github.io/ in Ukraine.
type Dictionary struct {
	XMLName                  xml.Name `xml:"Dictionary"`
	Text                     string   `xml:"chardata"`
	GeneratedByEngineVersion string   `xml:"GeneratedByEngineVersion,attr"`
	RequiredEngineVersion    string   `xml:"RequiredEngineVersion,attr"`
	Encrypted                string   `xml:"Encrypted,attr"`
	Encoding                 string   `xml:"IsUTF16,attr"`
	Format                   string   `xml:"Format,attr"`
	Stripkey                 string   `xml:"Stripkey,attr"`
	CreationDate             string   `xml:"creationDate,attr"`
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
