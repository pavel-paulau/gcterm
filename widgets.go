package main

import (
	"github.com/gizak/termui"
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
