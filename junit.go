package main

import (
	"encoding/xml"
	"io"
	"os"
	"path"

	"github.com/bmatcuk/doublestar"
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

func addFileTimesFromIOReader(fileTimes map[string]float64, reader io.Reader) {
	junitXML := loadJUnitXML(reader)
	for _, testCase := range junitXML.TestCases {
		filePath := path.Clean(testCase.File)
		fileTimes[filePath] += testCase.Time
	}
}

func getFileTimesFromJUnitXML(fileTimes map[string]float64) {
	if junitXMLPath != "" {
		filenames, err := doublestar.Glob(junitXMLPath)
		if err != nil {
			fatalMsg("failed to match jUnit filename pattern: %v", err)
		}
		for _, junitFilename := range filenames {
			file, err := os.Open(junitFilename)
			if err != nil {
				fatalMsg("failed to open junit xml: %v\n", err)
			}
			printMsg("using test times from JUnit report %s\n", junitXMLPath)
			addFileTimesFromIOReader(fileTimes, file)
			file.Close()
		}
	} else {
		printMsg("using test times from JUnit report at stdin\n")
		addFileTimesFromIOReader(fileTimes, os.Stdin)
	}
}
