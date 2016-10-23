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

func newPar(label string) *termui.Par {
	par := termui.NewPar("")

	par.BorderLabel = label
	par.BorderFg = termui.ColorWhite
	par.BorderLabelFg = termui.ColorWhite
	par.Height = barHeight
	par.PaddingLeft = 10

	return par
}

func newGauge(label string) *termui.Gauge {
	gauge := termui.NewGauge()

	gauge.BarColor = termui.ColorRed
	gauge.Border = true
	gauge.BorderLabel = label
	gauge.BorderLabelFg = termui.ColorWhite
	gauge.Height = barHeight
	gauge.LabelAlign = termui.AlignCenter
	gauge.PaddingLeft = 1
	gauge.PaddingRight = 1
	gauge.Percent = 0
	gauge.PercentColorHighlighted = termui.AttrBold

	return gauge
}

func newLineChart(label string) *termui.LineChart {
	lineChart := termui.NewLineChart()

	lineChart.AxesColor = termui.ColorWhite
	lineChart.BorderLabel = label
	lineChart.BorderLabelFg = termui.ColorWhite
	lineChart.Data = make([]float64, width)
	lineChart.DataLabels = make([]string, width)
	lineChart.Height = height
	lineChart.LineColor = termui.ColorGreen | termui.AttrBold
	lineChart.Mode = "dot"

	return lineChart
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
