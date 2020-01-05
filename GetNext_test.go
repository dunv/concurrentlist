package concurrentList

import (
	"testing"
	"time"
)

// This will get stuck in a deadlock, if it fails
func TestGetNext(t *testing.T) {
	list := NewConcurrentList()
	insertItems := []map[int]bool{}
	verifyItems := []map[int]bool{}

	totalProducer := 10
	totalItemsPerProducer := 10000
	totalConsumer := 1000
	bufferSize := totalProducer * totalItemsPerProducer

	// Create fixture
	readChannel := make(chan []int, bufferSize)
	for i := 0; i < totalProducer; i++ {
		insertItems = append(insertItems, map[int]bool{})
		verifyItems = append(verifyItems, map[int]bool{})
		for j := 0; j < totalItemsPerProducer; j++ {
			insertItems[i][j] = false
			verifyItems[i][j] = false
		}
	}

	// Create consumers
	for i := 0; i < totalConsumer; i++ {
		go consumer(list, &readChannel, t)
	}

	// Create producers
	for i := 0; i < totalProducer; i++ {
		go producer(insertItems[i], i, list)
	}

	// Validate
	for item := range readChannel {
		verifyItems[item[0]][item[1]] = true
		if verify(verifyItems) {
			return
		}
	}

}

func verify(verifyItems []map[int]bool) bool {
	for producerKey := range verifyItems {
		for _, itemValue := range verifyItems[producerKey] {
			if !itemValue {
				return false
			}
		}
	}
	return true
}

func consumer(list *ConcurrentList, readChannel *chan []int, t *testing.T) {
	timeSpentWaiting := int64(0)
	for {
		tmp := time.Now()
		item, err := list.GetNext()
		dur := time.Since(tmp)
		timeSpentWaiting += dur.Nanoseconds()
		if err != nil {
			t.Error("error", err)
			continue
		}

		parsed, ok := item.([]int)
		if ok {
			*readChannel <- parsed
		} else {
			t.Errorf("received unexpected item %v", item)
		}
	}
}

func producer(insertItems map[int]bool, producerIndex int, list *ConcurrentList) {
	for index := range insertItems {
		tmp1 := producerIndex
		tmp2 := index
		list.Append([]int{tmp1, tmp2})
	}
}