package main

import (
	"fmt"
	"sync"
	"time"
)

const iterations = 10000
const lineLength = 80
const controlTestName = "control test"
const lockingTestName = "test with locking"
const noLockingTestName = "test without locking"

type counter struct {
	mutex sync.Mutex
	x     int
}

func main() {
	timeWithLocking := measureTime(true, false, lockingTestName)
	timeWithoutLocking := measureTime(false, false, noLockingTestName)
	controlTime := measureTime(false, true, controlTestName)

	compareTimes(timeWithLocking, controlTime, lockingTestName, controlTestName)
	compareTimes(timeWithoutLocking, controlTime, noLockingTestName, controlTestName)
	compareTimes(timeWithLocking, timeWithoutLocking, lockingTestName, noLockingTestName)
}

func measureTime(lock bool, singleThread bool, header string) int64 {
	var waitGroup sync.WaitGroup
	counter := counter{x: 0}
	timeStart := time.Now().UnixNano()
	printLine(lineLength)
	fmt.Println(header)
	for i := 0; iterations > i; i++ {
		if singleThread {
			checkIfPrime(i)
			counter.x++
			continue
		}
		waitGroup.Add(1)
		go work(i, &waitGroup, &counter, lock)
	}
	if !singleThread {
		waitGroup.Wait()
	}
	timeEnd := time.Now().UnixNano()
	difference := timeEnd - timeStart
	fmt.Printf("Start: %d End: %d Diference: %d\n", timeStart, timeEnd, difference)
	return difference
}

func work(number int, wg *sync.WaitGroup, counter *counter, lock bool) {
	defer wg.Done()
	if lock {
		counter.mutex.Lock()
	}
	checkIfPrime(number)
	counter.x++
	if lock {
		counter.mutex.Unlock()
	}
}

func printLine(length int) {
	fmt.Print("\n")
	for i := 0; length > i; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n\n")
}

//this function isn't optimized on purpose
func checkIfPrime(number int) bool {
	dividers := 0
	for i := 2; i < number; i++ {
		if number%i == 0 {
			dividers++
		}
	}
	return dividers == 0
}

func compareTimes(timeOne int64, timeTwo int64, nameOne string, nameTwo string) {
	printLine(lineLength)
	fmt.Printf("Compare %s with %s:\n", nameOne, nameTwo)
	if timeOne > timeTwo {
		difference := timeOne - timeTwo
		percentage := float64(difference*100) / float64(timeTwo)
		fmt.Printf("%s took %d nanoseconds (%f %%) more\n\n", nameOne, difference, percentage)
		return
	}
	difference := timeTwo - timeOne
	percentage := float64(difference*100) / float64(timeOne)
	fmt.Printf("%s took %d nanoseconds (%f %%) more\n\n", nameTwo, difference, percentage)
}
