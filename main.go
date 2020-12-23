package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	guuid "github.com/google/uuid"
)

/*
Input ...
*/
type Input struct {
	uuid string
}

/*
Output ...
*/
type Output struct {
	hash string
}

func getMD5Hash(text string) string {
	time.Sleep(2 * time.Second)
	hasher := md5.New()
	hasher.Write([]byte(text))
	md5Value := hex.EncodeToString(hasher.Sum(nil))
	return md5Value
}

func generateRandomUUID() string {
	id := guuid.New()
	generatedUUID := id.String()
	return generatedUUID
}

func generateInputs() []Input {
	myUUIDList := make([]Input, 0)
	for i := 0; i < 15; i++ {
		myUUID := generateRandomUUID()
		myData := Input{uuid: myUUID}
		myUUIDList = append(myUUIDList, myData)
	}
	return myUUIDList
}

func getMD5HashForAllUUID() []Output {
	jobs := make(chan Input, 1)
	results := make(chan Output, 1)
	resultList := make([]Output, 0)
	maxNumberOfWorkers := 3

	uuidList := generateInputs()

	fmt.Printf("\nDone generating UUID List of Inputs.")

	waitChannel := make(chan struct{})

	for w := 1; w <= maxNumberOfWorkers; w++ {
		go hashWorker(jobs, results)
	}

	go func() {
		for _, uuid := range uuidList {
			jobs <- uuid
		}
		close(jobs)
	}()

	go func() {
		for i := 0; i < len(uuidList); i++ {
			result := <-results
			resultList = append(resultList, result)
		}
		close(waitChannel)
	}()

	<-waitChannel

	return resultList
}

func hashWorker(jobs <-chan Input, results chan<- Output) {
	counter := 1
	for {
		job, ok := <-jobs
		if !ok {
			break
		}
		hash := getMD5Hash(job.uuid)
		fmt.Printf("\n(%v) : hashWorker : job => %v , hash => %v", counter, job, hash)
		counter++
		output := Output{hash: hash}
		results <- output
	}
}

func main() {
	hashList := getMD5HashForAllUUID()
	fmt.Printf("\nTotal Length of output : %v", len(hashList))
	for _, data := range hashList {
		fmt.Printf("\nHash of data : %v", data.hash)
	}
}
