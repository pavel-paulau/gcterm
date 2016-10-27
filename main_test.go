package main

import (
	"testing"

	"github.com/gizak/termui"
)

func TestLayout(t *testing.T) {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	buildLayout() // Lack of crashes is sufficient
}
