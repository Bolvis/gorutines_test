package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type counter struct {
	mutex sync.Mutex
	x     int
}

const TimeFormat = "0.000000"
const BitSize = 64
const Iterations = 10000

func main() {

	timeOne := measureTime(true)

	for i := 0; 100 > i; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")

	timeTwo := measureTime(false)

	controlTime := singleThreadControlSumTime()

	fmt.Printf("Control time is %f\n", controlTime)

	if timeOne > timeTwo {
		fmt.Printf("\nWith locking took %f longer", timeOne-timeTwo)
		return
	}
	fmt.Printf("\nWithout locking took %f longer", timeTwo-timeOne)

}

func measureTime(lock bool) float64 {
	var waitGroup sync.WaitGroup
	counter := counter{x: 0}
	timeStart, _ := strconv.ParseFloat(time.Now().Format(TimeFormat), BitSize)
	for i := 0; Iterations > i; i++ {
		waitGroup.Add(1)
		go work(i, &waitGroup, &counter, lock)
	}
	waitGroup.Wait()
	timeEnd, _ := strconv.ParseFloat(time.Now().Format(TimeFormat), BitSize)
	fmt.Println("\nWork done\n")
	difference := timeEnd - timeStart
	fmt.Printf("Start: %f End: %f Diference: %f\n", timeStart, timeEnd, difference)
	return difference
}

func work(number int, wg *sync.WaitGroup, counter *counter, lock bool) {
	defer wg.Done()
	if lock {
		counter.mutex.Lock()
	}
	fmt.Printf("iteration: %d counter: %d\n", number, counter.x)
	counter.x++
	if lock {
		counter.mutex.Unlock()
	}
}

func singleThreadControlSumTime() float64 {
	timeStart, _ := strconv.ParseFloat(time.Now().Format(TimeFormat), BitSize)
	counter := 0
	for i := 0; Iterations > i; i++ {
		fmt.Printf("iteration: %d counter: %d\n", i, counter)
		counter++
	}
	timeEnd, _ := strconv.ParseFloat(time.Now().Format(TimeFormat), BitSize)
	fmt.Println("\nWork done\n")
	difference := timeEnd - timeStart
	fmt.Printf("Start: %f End: %f Diference: %f\n", timeStart, timeEnd, difference)
	return difference
}
