package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/gizak/termui"
)

var (
	mu                 sync.Mutex
	gcCounter, stwTime int
	gcPercent          *termui.Gauge
	gcRate             *termui.Par
	liveHeap, goalHeap *termui.LineChart
	wallTime, cpuTime  *termui.BarChart
)

func refreshGCSummary() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		gcRate.Text = strconv.Itoa(gcCounter / 5)
		gcPercent.Percent = 100 * stwTime / 5e6

		termui.Render(termui.Body)

		mu.Lock()
		gcCounter = 0
		stwTime = 0
		mu.Unlock()
	}
}

func refreshGraphs(data gcInfo) {
	liveHeap.Data = append(liveHeap.Data[1:], data.size.live)
	goalHeap.Data = append(goalHeap.Data[1:], data.size.goal)

	wallTime.Data = []int{
		data.wallTime.sweepTermination,
		data.wallTime.markAndSwap,
		data.wallTime.markTermination,
	}
	cpuTime.Data = []int{
		data.cpuTime.sweepTermination,
		data.cpuTime.markAndSwap,
		data.cpuTime.markTermination,
	}

	mu.Lock()
	gcCounter++
	stwTime += data.wallTime.sweepTermination
	stwTime += data.wallTime.markTermination
	mu.Unlock()
}

func sendEvents() {
	for data := range readStdin() {
		termui.SendCustomEvt("/feed", data)
	}
}

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	gcPercent = newGauge("Percentage of Time Spent in GC")
	gcRate = newPar("GC Events per Second")

	liveHeap = newLineChart("Live heap size, MB")
	goalHeap = newLineChart("Goal heap size, MB")

	wallTime = newBarChart("Wall-clock time, us",
		[]string{"STW Sweep Termination", "Concurrent Mark & Swap", "STW Mark Termination"})
	cpuTime = newBarChart("CPU time, us",
		[]string{"STW Sweep Termination", "Concurrent Mark & Swap", "STW Mark Termination"})

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(10, 0, gcPercent),
			termui.NewCol(2, 0, gcRate)),
		termui.NewRow(
			termui.NewCol(12, 0, liveHeap)),
		termui.NewRow(
			termui.NewCol(12, 0, goalHeap)),
		termui.NewRow(
			termui.NewCol(6, 0, wallTime),
			termui.NewCol(6, 0, cpuTime)))

	termui.Body.Align()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/feed", func(e termui.Event) {
		if data, ok := e.Data.(gcInfo); ok {
			refreshGraphs(data)

			termui.Render(termui.Body)
		}
	})

	termui.Render(termui.Body)

	go sendEvents()

	go refreshGCSummary()

	termui.Loop()
}
