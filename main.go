package main

import (
	"golang.org/x/crypto/ssh/terminal"

	"github.com/gizak/termui"
)

var (
	width, height int
)

func init() {
	termWidth, termHeight, err := terminal.GetSize(0)
	if err != nil {
		panic(err)
	}
	width = termWidth
	height = (termHeight - gaugeHeight) / 3
}

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	gcPercent := newGauge("Percentage of Time Spent in GC")

	gcRate := newPar("GC Events per Minute")

	liveHeap := newLineChart("Live heap size, MB")
	goalHeap := newLineChart("Goal heap size, MB")

	wallTime := newBarChart("Wall-clock time, us",
		[]string{"Sweep Termination", "Mark & Swap", "Mark Termination"}, false)
	cpuTime := newBarChart("CPU time, us",
		[]string{"Assist", "Background GC", "Idle GC"}, false)

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

	termui.Handle("/timer/1s", func(e termui.Event) {
		termui.Render(termui.Body)
	})

	termui.Loop()
}
