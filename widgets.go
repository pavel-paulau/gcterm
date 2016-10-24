package main

import (
	"github.com/gizak/termui"
)

const (
	gaugeHeight = 3
	yAxisWidth  = 15
)

func newPar(label string) *termui.Par {
	par := termui.NewPar("N/A")

	par.BorderLabel = label
	par.BorderFg = termui.ColorWhite
	par.BorderLabelFg = termui.ColorYellow
	par.Height = gaugeHeight
	par.PaddingLeft = 10

	return par
}

func newGauge(label string) *termui.Gauge {
	gauge := termui.NewGauge()

	gauge.BarColor = termui.ColorRed
	gauge.Border = true
	gauge.BorderLabel = label
	gauge.BorderLabelFg = termui.ColorYellow
	gauge.Height = gaugeHeight
	gauge.LabelAlign = termui.AlignCenter
	gauge.PaddingLeft = 1
	gauge.PaddingRight = 1
	gauge.Percent = 0
	gauge.PercentColorHighlighted = termui.AttrBold

	return gauge
}

func newLineChart(label string) *termui.LineChart {
	height := (termui.TermHeight() - gaugeHeight) / 3
	width := termui.TermWidth() - yAxisWidth

	lineChart := termui.NewLineChart()

	lineChart.AxesColor = termui.ColorWhite
	lineChart.BorderLabel = label
	lineChart.BorderLabelFg = termui.ColorYellow
	lineChart.Data = make([]float64, width)
	lineChart.DataLabels = make([]string, width)
	lineChart.Height = height
	lineChart.LineColor = termui.ColorYellow | termui.AttrBold
	lineChart.Mode = "dot"

	return lineChart
}

func newBarChart(label string, dataLabels []string) *termui.BarChart {
	height := (termui.TermHeight() - gaugeHeight) / 3
	width := termui.TermWidth() / 2

	bc := termui.NewBarChart()

	bc.BarColor = termui.ColorBlue
	bc.BarWidth = (width/len(dataLabels) - len(dataLabels))
	bc.BorderLabel = label
	bc.BorderLabelFg = termui.ColorYellow
	bc.Data = make([]int, len(dataLabels))
	bc.DataLabels = dataLabels
	bc.Height = height
	bc.NumColor = termui.ColorWhite | termui.AttrBold
	bc.PaddingLeft = 1
	bc.TextColor = termui.ColorWhite

	return bc
}
