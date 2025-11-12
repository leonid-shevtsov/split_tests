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

// loadJUnitTimingsFromGlob loads test timings from JUnit XML files matching a glob pattern
func loadJUnitTimingsFromGlob(globPattern string) map[string]float64 {
	fileTimes := make(map[string]float64)

	if globPattern == "" {
		return fileTimes
	}

	filenames, err := doublestar.Glob(globPattern)
	if err != nil {
		fatalMsg("failed to match jUnit filename pattern: %v", err)
	}

	if len(filenames) == 0 {
		printMsg("warning: no files matched pattern %s\n", globPattern)
		return fileTimes
	}

	for _, junitFilename := range filenames {
		file, err := os.Open(junitFilename)
		if err != nil {
			fatalMsg("failed to open junit xml: %v\n", err)
		}
		printMsg("loaded test times from %s\n", junitFilename)
		addFileTimesFromIOReader(fileTimes, file)
		file.Close()
	}

	return fileTimes
}

func getFileTimesFromJUnitXML(fileTimes map[string]float64) {
	if junitXMLPath != "" {
		loadedTimes := loadJUnitTimingsFromGlob(junitXMLPath)
		for file, time := range loadedTimes {
			fileTimes[file] += time
		}
	} else {
		printMsg("using test times from JUnit report at stdin\n")
		addFileTimesFromIOReader(fileTimes, os.Stdin)
	}
}
