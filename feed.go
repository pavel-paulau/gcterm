package main

import (
	"bufio"
	"errors"
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
func parseClock(s string) (gcTimings, error) {
	timings := strings.Split(s, "+")

	if len(timings) != 3 {
		return gcTimings{}, errors.New("bad wall-clock timings")
	}

	sweepTermination, err := strconv.ParseFloat(timings[0], 64)
	if err != nil {
		return gcTimings{}, err
	}
	markAndSwap, err := strconv.ParseFloat(timings[1], 64)
	if err != nil {
		return gcTimings{}, err
	}
	markTermination, err := strconv.ParseFloat(timings[2], 64)
	if err != nil {
		return gcTimings{}, err
	}

	return gcTimings{
		sweepTermination: int(sweepTermination * 1e3),
		markAndSwap:      int(markAndSwap * 1e3),
		markTermination:  int(markTermination * 1e3),
	}, nil
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
func parseCPU(s string) (gcTimings, error) {
	timings := strings.Split(s, "+")
	if len(timings) != 3 {
		return gcTimings{}, errors.New("bad CPU timings")
	}

	markAndSwapTimings := strings.Split(timings[1], "/")
	if len(markAndSwapTimings) != 3 {
		return gcTimings{}, errors.New("bad mark and swap timings")
	}

	assist, err := strconv.ParseFloat(markAndSwapTimings[0], 64)
	if err != nil {
		return gcTimings{}, err
	}
	background, err := strconv.ParseFloat(markAndSwapTimings[1], 64)
	if err != nil {
		return gcTimings{}, err
	}
	idle, err := strconv.ParseFloat(markAndSwapTimings[2], 64)
	if err != nil {
		return gcTimings{}, err
	}
	markAndSwap := assist + background + idle

	sweepTermination, err := strconv.ParseFloat(timings[0], 64)
	if err != nil {
		return gcTimings{}, err
	}
	markTermination, err := strconv.ParseFloat(timings[2], 64)
	if err != nil {
		return gcTimings{}, err
	}

	return gcTimings{
		sweepTermination: int(sweepTermination * 1e3),
		markAndSwap:      int(markAndSwap * 1e3),
		markTermination:  int(markTermination * 1e3),
	}, nil
}

// parseLive parses gctrace output and extracts live heap size.
//
// 	#->#-># MB
//
func parseLive(s string) (float64, error) {
	sizes := strings.Split(s, "->")

	size, err := strconv.ParseInt(sizes[2], 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(size), nil
}

// parseGoal parses gctrace output and extracts goal heap size.
//
// 	# MB goal
//
func parseGoal(s string) (float64, error) {
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(size), nil
}

// readStdin reads gctrace lines from standard input and emits information about
// GC events.
//
// gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
//
func readStdin() <-chan gcInfo {
	ch := make(chan gcInfo)

	go func() {
		defer close(ch)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fields := strings.Split(scanner.Text(), " ")
			if len(fields) != 17 {
				continue
			}

			info := gcInfo{}

			wallTime, err := parseClock(fields[4])
			if err != nil {
				continue
			}
			info.wallTime = wallTime

			cpuTime, err := parseCPU(fields[7])
			if err != nil {
				continue
			}
			info.cpuTime = cpuTime

			live, err := parseLive(fields[10])
			if err != nil {
				continue
			}
			goal, err := parseGoal(fields[12])
			if err != nil {
				continue
			}
			info.size = heapSize{
				live: live,
				goal: goal,
			}

			ch <- info
		}
	}()

	return ch
}
