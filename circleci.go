package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

func getCircleAPIJSON(url string, destination interface{}) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fatalMsg("error calling CircleCI API at %v: %v", url, err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(destination)
	if err != nil {
		fatalMsg("error parsing CircleCI JSON at %v: %v", url, err)
	}
}

type circleCIBranchList []struct {
	BuildNum int `json:"build_num"`
}

type circleCITestResults struct {
	Tests []struct {
		File    string  `json:"file"`
		RunTime float64 `json:"run_time"`
	} `json:"tests"`
}

func circleCIAPIURL() string {
	return fmt.Sprintf("https://circleci.com/api/v1.1/project/%s", circleCIProjectPrefix)
}

func getCircleCIBranchBuilds(branchName string) circleCIBranchList {
	buildsURL := fmt.Sprintf("%s/tree/%s?filter=successful&circle-token=%s", circleCIAPIURL(), branchName, circleCIAPIKey)
	var branchList circleCIBranchList
	getCircleAPIJSON(buildsURL, &branchList)
	return branchList
}

func getCircleCITestResults(buildNum int) circleCITestResults {
	testResultsURL := fmt.Sprintf("%s/%d/tests?circle-token=%s", circleCIAPIURL(), buildNum, circleCIAPIKey)
	var testResults circleCITestResults
	getCircleAPIJSON(testResultsURL, &testResults)
	return testResults
}

func getFileTimesFromCircleCI(fileTimes map[string]float64) {
	builds := getCircleCIBranchBuilds(circleCIBranchName)
	if len(builds) == 0 {
		builds = getCircleCIBranchBuilds("master")
	}
	if len(builds) > 0 {
		buildNum := builds[0].BuildNum
		printMsg("using test timings from CircleCI build %d\n", buildNum)

		testResults := getCircleCITestResults(buildNum)

		for _, test := range testResults.Tests {
			filePath := path.Clean(test.File)
			fileTimes[filePath] += test.RunTime
		}
	}
}
