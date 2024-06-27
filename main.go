package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar"
)

var useCircleCI bool
var useJUnitXML bool
var useLineCount bool
var junitXMLPath string
var testFilePatterns []string
var excludeFilePattern = ""
var circleCIProjectPrefix = ""
var circleCIBranchName string
var splitIndex int
var splitTotal int
var circleCIAPIKey string

func printMsg(msg string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func fatalMsg(msg string, args ...interface{}) {
	printMsg(msg, args...)
	os.Exit(1)
}

func removeDeletedFiles(fileTimes map[string]float64, currentFileSet map[string]bool) {
	for file := range fileTimes {
		if !currentFileSet[file] {
			delete(fileTimes, file)
		}
	}
}

func addNewFiles(fileTimes map[string]float64, currentFileSet map[string]bool) {
	averageFileTime := 0.0
	if len(fileTimes) > 0 {
		for _, time := range fileTimes {
			averageFileTime += time
		}
		averageFileTime /= float64(len(fileTimes))
	} else {
		averageFileTime = 1.0
	}

	for file := range currentFileSet {
		if _, isSet := fileTimes[file]; isSet {
			continue
		}
		if useCircleCI || useJUnitXML {
			printMsg("missing file time for %s\n", file)
		}
		fileTimes[file] = averageFileTime
	}
}

func parseFlags() {
	var testFilePattern string
	flag.StringVar(&testFilePattern, "glob", "spec/**/*_spec.rb", "Comma-separated glob patterns to find test files. Make sure to single-quote to avoid shell expansion.")
	flag.StringVar(&excludeFilePattern, "exclude-glob", "", "Glob pattern to exclude test files. Make sure to single-quote.")

	flag.IntVar(&splitIndex, "split-index", -1, "This test container's index (or set CIRCLE_NODE_INDEX)")
	flag.IntVar(&splitTotal, "split-total", -1, "Total number of containers (or set CIRCLE_NODE_TOTAL)")

	flag.StringVar(&circleCIAPIKey, "circleci-key", "", "CircleCI API key (or set CIRCLECI_API_KEY environment variable) - required to use CircleCI")
	flag.StringVar(&circleCIProjectPrefix, "circleci-project", "", "CircleCI project name (e.g. github/leonid-shevtsov/split_tests) - required to use CircleCI")
	flag.StringVar(&circleCIBranchName, "circleci-branch", "", "Current branch for CircleCI (or set CIRCLE_BRANCH) - required to use CircleCI")

	flag.BoolVar(&useJUnitXML, "junit", false, "Use a JUnit XML report for test times")
	flag.StringVar(&junitXMLPath, "junit-path", "", "Path to a JUnit XML report (leave empty to read from stdin; use glob pattern to load multiple files)")

	flag.BoolVar(&useLineCount, "line-count", false, "Use line count to estimate test times")

	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show this help text")

	flag.Parse()

	testFilePatterns = strings.Split(testFilePattern, ",")

	var err error
	if circleCIAPIKey == "" {
		circleCIAPIKey = os.Getenv("CIRCLECI_API_KEY")
	}
	if circleCIBranchName == "" {
		circleCIBranchName = os.Getenv("CIRCLE_BRANCH")
	}
	if splitTotal == -1 {
		splitTotal, err = strconv.Atoi(os.Getenv("CIRCLE_NODE_TOTAL"))
		if err != nil {
			splitIndex = -1
		}
	}
	if splitIndex == -1 {
		splitIndex, err = strconv.Atoi(os.Getenv("CIRCLE_NODE_INDEX"))
		if err != nil {
			splitIndex = -1
		}
	}

	useCircleCI = circleCIAPIKey != ""

	if showHelp {
		printMsg("Splits test files into containers of even duration\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if useCircleCI && (circleCIProjectPrefix == "" || circleCIBranchName == "") {
		fatalMsg("Incomplete CircleCI configuration (set -circleci-key, -circleci-project, and -circleci-branch\n")
	}
	if splitTotal == 0 || splitIndex < 0 || splitIndex > splitTotal {
		fatalMsg("-split-index and -split-total (and environment variables) are missing or invalid\n")
	}
}

func main() {
	parseFlags()

	// We are not using filepath.Glob,
	// because it doesn't support '**' (to match all files in all nested directories)
	currentFileSet := make(map[string]bool)
	for _, pattern := range testFilePatterns {
		currentFiles, err := doublestar.Glob(pattern)
		if err != nil {
			printMsg("failed to enumerate current file set for pattern %s: %v", pattern, err)
			os.Exit(1)
		}
		for _, file := range currentFiles {
			currentFileSet[file] = true
		}
	}

	if excludeFilePattern != "" {
		excludedFiles, err := doublestar.Glob(excludeFilePattern)
		if err != nil {
			printMsg("failed to enumerate excluded file set: %v", err)
			os.Exit(1)
		}
		for _, file := range excludedFiles {
			delete(currentFileSet, file)
		}
	}

	fileTimes := make(map[string]float64)
	if useLineCount {
		estimateFileTimesByLineCount(currentFileSet, fileTimes)
	} else if useJUnitXML {
		getFileTimesFromJUnitXML(fileTimes)
	} else if useCircleCI {
		getFileTimesFromCircleCI(fileTimes)
	}

	removeDeletedFiles(fileTimes, currentFileSet)
	addNewFiles(fileTimes, currentFileSet)

	buckets, bucketTimes := splitFiles(fileTimes, splitTotal)
	if useCircleCI || useJUnitXML {
		printMsg("expected test time: %0.1fs\n", bucketTimes[splitIndex])
	}

	fmt.Println(strings.Join(buckets[splitIndex], " "))
}
