package main

import (
	"encoding/xml"
	"io"
	"os"
	"path"
)

type junitXML struct {
	TestCases []struct {
		File string  `xml:"file,attr"`
		Time float64 `xml:"time,attr"`
	} `xml:"testcase"`
}

func loadJUnitXML(reader io.Reader) *junitXML {
	var junitXML junitXML

	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&junitXML)
	if err != nil {
		fatalMsg("failed to parse junit xml: %v\n", err)
	}

	return &junitXML
}

func getFileTimesFromJUnitXML(fileTimes map[string]float64) {
	var source io.Reader
	if junitXMLPath != "" {
		file, err := os.Open(junitXMLPath)
		if err != nil {
			fatalMsg("failed to open junit xml: %v\n", err)
		}
		defer file.Close()
		printMsg("using test times from JUnit report %s\n", junitXMLPath)
		source = file
	} else {
		printMsg("using test times from JUnit report at stdin\n")
		source = os.Stdin
	}
	junitXML := loadJUnitXML(source)
	for _, testCase := range junitXML.TestCases {
		filePath := path.Clean(testCase.File)
		fileTimes[filePath] += testCase.Time
	}
}
