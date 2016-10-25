package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/gizak/termui"
)

const (
	calcInterval      = 2
	renderingInterval = 250 * time.Millisecond
)

var (
	mu                 sync.Mutex
	gcCounter, stwTime int
	gcPercent          *termui.Gauge
	gcRate             *termui.Par
	liveHeap, goalHeap *termui.LineChart
	wallTime, cpuTime  *termui.BarChart
)

func render() {
	ticker := time.NewTicker(renderingInterval)

	for range ticker.C {
		termui.Render(termui.Body)
	}
}

func updateSummary() {
	ticker := time.NewTicker(calcInterval * time.Second)

	for range ticker.C {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			gcRate.Text = strconv.Itoa(gcCounter / calcInterval)
			gcPercent.Percent = 100 * stwTime / calcInterval / 1e6

			gcCounter = 0
			stwTime = 0
		}()
	}
}

func updateCharts(data gcInfo) {
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

func initWidgets() {
	gcPercent = newGauge("Percentage of Time Spent in STW GC")
	gcRate = newPar("GC Events per Second")

	liveHeap = newLineChart("Live heap size, MB")
	goalHeap = newLineChart("Goal heap size, MB")

	phases := []string{
		"STW Sweep Termination",
		"Concurrent Mark & Swap",
		"STW Mark Termination",
	}
	wallTime = newBarChart("Wall-clock time, us", phases)
	cpuTime = newBarChart("CPU time, us", phases)
}

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	initWidgets()

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
			updateCharts(data)
		}
	})

	go render()

	go sendEvents()

	go updateSummary()

	termui.Loop()
}
