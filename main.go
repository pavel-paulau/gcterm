package main

import (
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gizak/termui"
)

const (
	barHeight  = 3
	maxTicks   = 4
	timeFmt    = "15:04:05"
	xAxisWidth = 15
	xTime      = 30 * time.Second
)

var (
	width, height int
)

func init() {
	var err error
	width, height, err = terminal.GetSize(0)
	if err != nil {
		panic(err)
	}
	width -= xAxisWidth
	height = (height - barHeight) / 2
}

func addTimeLabels(lc *termui.LineChart) {
	t := time.Now()

	for tick := 0; tick < maxTicks; tick++ {
		tickIdx := tick * width / maxTicks
		offset := time.Duration(maxTicks-tick) * (xTime / maxTicks)
		lc.DataLabels[tickIdx] = t.Add(-offset).Format(timeFmt)
	}
}

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	gcRate := newPar("GC events per minute")
	gcPercent := newGauge("Percentage of time spent in GC since program start")

	liveHeap := newLineChart("Live heap size, MB")
	goalHeap := newLineChart("Goal heap size, MB")

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(10, 0, gcPercent),
			termui.NewCol(2, 0, gcRate)),
		termui.NewRow(
			termui.NewCol(12, 0, liveHeap)),
		termui.NewRow(
			termui.NewCol(12, 0, goalHeap)))

	termui.Body.Align()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		addTimeLabels(liveHeap)
		addTimeLabels(goalHeap)

		termui.Render(termui.Body)
	})

	termui.Loop()
}
