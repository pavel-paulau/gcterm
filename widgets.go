package main

import (
	"github.com/gizak/termui"
)

const (
	gaugeHeight = 3
	xAxisWidth  = 15
)

func newPar(label string) *termui.Par {
	par := termui.NewPar("")

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
	lineChart := termui.NewLineChart()

	lineChart.AxesColor = termui.ColorWhite
	lineChart.BorderLabel = label
	lineChart.BorderLabelFg = termui.ColorYellow
	lineChart.Data = make([]float64, width-xAxisWidth)
	lineChart.DataLabels = make([]string, width-xAxisWidth)
	lineChart.Height = height
	lineChart.LineColor = termui.ColorYellow | termui.AttrBold
	lineChart.Mode = "dot"

	return lineChart
}

func newBarChart(label string, dataLabels []string, fullSize bool) *termui.BarChart {
	bc := termui.NewBarChart()

	bc.BarColor = termui.ColorBlue
	bc.BarGap = 2
	bc.BorderLabel = label
	bc.BorderLabelFg = termui.ColorYellow
	bc.Data = make([]int, len(dataLabels))
	bc.DataLabels = dataLabels
	bc.Height = height
	bc.NumColor = termui.ColorWhite | termui.AttrBold
	bc.PaddingLeft = 1
	bc.TextColor = termui.ColorWhite

	if fullSize {
		bc.BarWidth = (width - len(dataLabels)*bc.BarGap - len(dataLabels)) / len(dataLabels)
	} else {
		bc.BarWidth = (width/2 - len(dataLabels)*bc.BarGap - len(dataLabels)) / len(dataLabels)
	}
	return bc
}
