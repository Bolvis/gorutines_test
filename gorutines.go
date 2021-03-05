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

func main() {

	fmt.Println("Locking")
	timeOne := measureTime(true)
	printLine(100)
	fmt.Println("No locking")
	timeTwo := measureTime(false)
	printLine(100)
	fmt.Println("Control")
	controlTime := singleThreadControlSumTime()

	fmt.Printf("\nControl time is %d\n", controlTime)

	if timeOne > timeTwo {
		difference := timeOne - timeTwo
		percentage := float64(difference*100) / float64(timeOne)
		fmt.Printf("\nWith locking took %d (%f %%) longer\n", difference, percentage)
		return
	}
	difference := timeTwo - timeOne
	percentage := float64(difference*100) / float64(timeTwo)
	fmt.Printf("\nWith locking took %d (%f %%) longer\n", difference, percentage)

}

func measureTime(lock bool) int64 {
	var waitGroup sync.WaitGroup
	counter := counter{x: 0}
	timeStart := time.Now().UnixNano()
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
	for i := 0; length > i; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
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
