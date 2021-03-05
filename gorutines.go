package main

import (
	"fmt"
	"sync"
	"time"
)

type counter struct {
	mutex sync.Mutex
	x     int
}

const Iterations = 10000
const LineLength = 80

func main() {

	timeOne := measureTime(true)
	timeTwo := measureTime(false)
	controlTime := singleThreadControlSumTime()

	printControlResults(timeOne, controlTime, "with locking")
	printControlResults(timeTwo, controlTime, "without locking")
	printResults(timeOne, timeTwo)

}

func measureTime(lock bool) int64 {
	var waitGroup sync.WaitGroup
	counter := counter{x: 0}
	timeStart := time.Now().UnixNano()
	printLine(LineLength)
	if lock {
		fmt.Println("With locking")
	} else {
		fmt.Println("Without locking")
	}
	for i := 0; Iterations > i; i++ {
		waitGroup.Add(1)
		go work(i, &waitGroup, &counter, lock)
	}
	waitGroup.Wait()
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
	//fmt.Printf("iteration: %d counter: %d\n", number, counter.x)
	checkIfPrime(number)
	counter.x++
	if lock {
		counter.mutex.Unlock()
	}
}

func singleThreadControlSumTime() int64 {
	timeStart := time.Now().UnixNano()
	counter := 0
	printLine(LineLength)
	fmt.Println("Control task")
	for i := 0; Iterations > i; i++ {
		//fmt.Printf("iteration: %d counter: %d\n", i, counter)
		checkIfPrime(i)
		counter++
	}
	timeEnd := time.Now().UnixNano()
	difference := timeEnd - timeStart
	fmt.Printf("Start: %d End: %d Diference: %d\n", timeStart, timeEnd, difference)
	return difference
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

func printResults(timeOne int64, timeTwo int64) {
	printLine(LineLength)
	if timeOne > timeTwo {
		difference := timeOne - timeTwo
		percentage := float64(difference*100) / float64(timeTwo)
		fmt.Printf("With locking took %d nanoseconds (%f %%) more\n\n", difference, percentage)
		return
	}
	difference := timeTwo - timeOne
	percentage := float64(difference*100) / float64(timeOne)
	fmt.Printf("With locking took %d nanoseconds (%f %%) more\n\n", difference, percentage)
}

func printControlResults(testTime int64, controlTime int64, testName string) {
	printLine(LineLength)
	fmt.Printf("Compare control task with \"%s\" task:\n", testName)
	if testTime > controlTime {
		difference := testTime - controlTime
		percentage := float64(difference*100) / float64(controlTime)
		fmt.Printf("Test %s took %d nanoseconds (%f %%) more\n", testName, difference, percentage)
		return
	}
	difference := controlTime - testTime
	percentage := float64(difference*100) / float64(testTime)
	fmt.Printf("Control task took %d nanoseconds (%f %%) more\n", difference, percentage)
}
