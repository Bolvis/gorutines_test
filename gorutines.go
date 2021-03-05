package main

import (
	"fmt"
	"sync"
	"time"
)

const Iterations = 10000
const LineLength = 80
const ControlTestName = "control test"
const LockingTestName = "test with locking"
const NoLockingTestName = "test without locking"

type counter struct {
	mutex sync.Mutex
	x     int
}

func main() {
	timeWithLocking := measureTime(true, false, LockingTestName)
	timeWithoutLocking := measureTime(false, false, NoLockingTestName)
	controlTime := measureTime(false, true, ControlTestName)

	compareTimes(timeWithLocking, controlTime, LockingTestName, ControlTestName)
	compareTimes(timeWithoutLocking, controlTime, NoLockingTestName, ControlTestName)
	compareTimes(timeWithLocking, timeWithoutLocking, LockingTestName, NoLockingTestName)
}

func measureTime(lock bool, singleThread bool, header string) int64 {
	var waitGroup sync.WaitGroup
	counter := counter{x: 0}
	timeStart := time.Now().UnixNano()
	printLine(LineLength)
	fmt.Println(header)
	for i := 0; Iterations > i; i++ {
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
	printLine(LineLength)
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
