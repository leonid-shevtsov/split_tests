package main

import "sort"

func splitFiles(biases []float64, fileTimesMap map[string]float64, splitTotal int) ([][]string, []float64) {
	buckets := make([][]string, splitTotal)
	bucketTimes := make([]float64, splitTotal)

	// Build a sorted list of files
	fileTimesList := make(fileTimesList, len(fileTimesMap))
	for file, time := range fileTimesMap {
		fileTimesList = append(fileTimesList, fileTimesListItem{file, time})
	}
	sort.Sort(fileTimesList)

	for _, file := range fileTimesList {
		// find bucket with min weight
		minBucket := 0
		for bucket := 1; bucket < splitTotal; bucket++ {
			if bucketTimes[bucket]+biases[bucket] < bucketTimes[minBucket]+biases[minBucket] {
				minBucket = bucket
			}
		}
		// add file to bucket
		buckets[minBucket] = append(buckets[minBucket], file.name)
		bucketTimes[minBucket] += file.time
	}

	return buckets, bucketTimes
}

type fileTimesListItem struct {
	name string
	time float64
}

type fileTimesList []fileTimesListItem

func (l fileTimesList) Len() int { return len(l) }

// Sorts by time descending, then by name ascending
// Sort by name is required for deterministic order across machines
func (l fileTimesList) Less(i, j int) bool {
	return l[i].time > l[j].time ||
		(l[i].time == l[j].time && l[i].name < l[j].name)
}

func (l fileTimesList) Swap(i, j int) {
	temp := l[i]
	l[i] = l[j]
	l[j] = temp
}
