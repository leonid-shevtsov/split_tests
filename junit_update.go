package main

import (
	"encoding/xml"
	"os"
	"strconv"
)

var junitUpdateOldGlob string
var junitUpdateNewGlob string
var junitUpdateOutPath string

const slidingWindowOldWeight = 0.9

// applySlidingWindow applies an exponential moving average to smooth out timing fluctuations
// Uses exponential moving average: oldWeight * old + (1 - oldWeight) * new
// This gives more weight to historical data while still incorporating recent changes
func applySlidingWindow(oldTime, newTime float64) float64 {
	return slidingWindowOldWeight*oldTime + (1-slidingWindowOldWeight)*newTime
}


// updateJUnitTimings merges old and new JUnit timings using a sliding window algorithm
func updateJUnitTimings() {
	if junitUpdateOldGlob == "" || junitUpdateNewGlob == "" || junitUpdateOutPath == "" {
		fatalMsg("junit-update requires -junit-update, -junit-new, and -junit-out flags\n")
	}

	// Load old timings
	oldTimings := loadJUnitTimingsFromGlob(junitUpdateOldGlob)
	printMsg("loaded %d test files from old timings\n", len(oldTimings))

	// Load new timings
	newTimings := loadJUnitTimingsFromGlob(junitUpdateNewGlob)
	printMsg("loaded %d test files from new timings\n", len(newTimings))

	// Merge timings using sliding window algorithm
	mergedTimings := make(map[string]float64)

	// Process all tests from new timings
	for file, newTime := range newTimings {
		oldTime, exists := oldTimings[file]
		if exists {
			// Test exists in both: apply sliding window
			mergedTimings[file] = applySlidingWindow(oldTime, newTime)
		} else {
			// Test not in old: use new timing
			mergedTimings[file] = newTime
		}
	}

	// Tests not in new are automatically excluded (not added to mergedTimings)

	printMsg("merged %d test files (removed tests not in new, used sliding window for existing tests)\n", len(mergedTimings))

	// Write output JUnit XML
	writeJUnitXML(mergedTimings, junitUpdateOutPath)
	printMsg("wrote updated timings to %s\n", junitUpdateOutPath)
}

// writeJUnitXML writes test timings to a JUnit XML file
func writeJUnitXML(timings map[string]float64, outputPath string) {
	// JUnit XML structure for writing (with testsuite root element)
	type testCase struct {
		File string `xml:"file,attr"`
		Time string `xml:"time,attr"`
	}

	type testSuite struct {
		XMLName   xml.Name   `xml:"testsuite"`
		Name      string     `xml:"name,attr"`
		Tests     int        `xml:"tests,attr"`
		TestCases []testCase `xml:"testcase"`
	}

	// Convert map to slice for consistent output
	testCases := make([]testCase, 0, len(timings))
	for file, time := range timings {
		// Format as decimal without scientific notation
		timeStr := strconv.FormatFloat(time, 'f', -1, 64)
		testCases = append(testCases, testCase{
			File: file,
			Time: timeStr,
		})
	}

	suite := testSuite{
		Name:      "rspec",
		Tests:     len(testCases),
		TestCases: testCases,
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		fatalMsg("failed to create output file: %v\n", err)
	}
	defer file.Close()

	// Write XML header
	file.WriteString(xml.Header)

	// Create encoder
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	// Encode XML
	err = encoder.Encode(suite)
	if err != nil {
		fatalMsg("failed to encode JUnit XML: %v\n", err)
	}

	file.WriteString("\n")
}
