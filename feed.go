package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type gcInfo struct {
	size     heapSize
	wallTime gcTimings
	cpuTime  gcTimings
}

type heapSize struct {
	goal, live float64
}

type gcTimings struct {
	sweepTermination, markAndSwap, markTermination int
}

// parseClock parses gctrace output and extracts times for the phases of the GC.
//
// 	#+#+# ms clock
//
// The phases are stop-the-world (STW) sweep termination, concurrent  mark and
// scan, and STW mark termination.
//
// All timings are converted to microseconds for compatibility with the bar
// charts.
func parseClock(s string) gcTimings {
	timings := strings.Split(s, "+")

	sweepTermination, err := strconv.ParseFloat(timings[0], 64)
	if err != nil {
		panic(err)
	}
	markAndSwap, err := strconv.ParseFloat(timings[1], 64)
	if err != nil {
		panic(err)
	}
	markTermination, err := strconv.ParseFloat(timings[2], 64)
	if err != nil {
		panic(err)
	}

	return gcTimings{
		sweepTermination: int(sweepTermination * 1e3),
		markAndSwap:      int(markAndSwap * 1e3),
		markTermination:  int(markTermination * 1e3),
	}
}

// parseCPU parses gctrace output and extracts wall-clock times for the phases
// of the GC.
//
// 	#+#/#/#+# ms cpu
//
// The phases are stop-the-world (STW) sweep termination, concurrent  mark and
// scan, and STW mark termination. The CPU times  for mark/scan are broken down
// in to assist time (GC performed in line with allocation), background GC time,
// and idle GC time.
//
// All timings are converted to microseconds for compatibility with the bar
// charts.
func parseCPU(s string) gcTimings {
	timings := strings.Split(s, "+")
	markAndSwapTimings := strings.Split(timings[1], "/")

	assist, err := strconv.ParseFloat(markAndSwapTimings[0], 64)
	if err != nil {
		panic(err)
	}
	background, err := strconv.ParseFloat(markAndSwapTimings[1], 64)
	if err != nil {
		panic(err)
	}
	idle, err := strconv.ParseFloat(markAndSwapTimings[2], 64)
	if err != nil {
		panic(err)
	}

	sweepTermination, err := strconv.ParseFloat(timings[0], 64)
	if err != nil {
		panic(err)
	}
	markAndSwap := assist + background + idle
	markTermination, err := strconv.ParseFloat(timings[2], 64)
	if err != nil {
		panic(err)
	}

	return gcTimings{
		sweepTermination: int(sweepTermination * 1e3),
		markAndSwap:      int(markAndSwap * 1e3),
		markTermination:  int(markTermination * 1e3),
	}
}

// parseLive parses gctrace output and extracts live heap size.
//
// 	#->#-># MB
//
func parseLive(s string) float64 {
	sizes := strings.Split(s, "->")

	size, err := strconv.ParseInt(sizes[2], 10, 64)
	if err != nil {
		panic(err)
	}

	return float64(size)
}

// parseGoal parses gctrace output and extracts goal heap size.
//
// 	# MB goal
//
func parseGoal(s string) float64 {
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return float64(size)
}

// readStdin reads gctrace lines from standard input and emits information about
// GC events.
//
// gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
//
func readStdin() <-chan gcInfo {
	data := make(chan gcInfo)

	go func() {
		defer close(data)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fields := strings.Split(scanner.Text(), " ")

			data <- gcInfo{
				size: heapSize{
					live: parseLive(fields[10]),
					goal: parseGoal(fields[12]),
				},
				wallTime: parseClock(fields[4]),
				cpuTime:  parseCPU(fields[7]),
			}
		}
	}()

	return data
}
