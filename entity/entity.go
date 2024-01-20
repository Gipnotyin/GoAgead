package entity

import "encoding/xml"

type RequestCurrency struct {
	Date  string `json:"date"`
	Value string `json:"val"`
}

type ResponseCurrency struct {
	Value string `xml:"Value"`
}

type XMLAnswer struct {
	XMLName xml.Name `xml:"ValCurs"`
	Record  []struct {
		Nominal   int    `xml:"Nominal"`
		Value     string `xml:"Value"`
		VunitRate string `xml:"VunitRate"`
	}
}
