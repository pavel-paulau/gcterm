package main

import (
	"golang.org/x/crypto/ssh/terminal"

	"github.com/gizak/termui"
)

const (
	barHeight   = 3
	minHeight   = 10
	xResolution = 200
)

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

func maxHeight() int {
	_, height, err := terminal.GetSize(0)
	if err != nil {
		return minHeight
	}
	return (height - barHeight) / 2
}

func newLineChart(label string) *termui.LineChart {
	lineChart := termui.NewLineChart()

	lineChart.AxesColor = termui.ColorWhite
	lineChart.BorderLabel = label
	lineChart.BorderLabelFg = termui.ColorWhite
	lineChart.Data = make([]float64, xResolution)
	lineChart.DataLabels = make([]string, xResolution)
	lineChart.Height = maxHeight()
	lineChart.LineColor = termui.ColorGreen | termui.AttrBold
	lineChart.Mode = "dot"

	return lineChart
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
		termui.Render(termui.Body)
	})

	termui.Loop()
}
