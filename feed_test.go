package main

import (
	"testing"
)

func TestParserGo16Go17(t *testing.T) {
	line := "gc 24 @0.252s 11%: 0.016+0.79+0.23 ms clock, 0.13+0.95/1.3/2.2+1.8 ms cpu, 8->9->2 MB, 12 MB goal, 8 P"

	i, err := extractInfo(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Wall-clock timings
	if i.wallTime.sweepTermination != 16 {
		t.Errorf("invalid wallTime.sweepTermination: got %v", i.wallTime.sweepTermination)
	}
	if i.wallTime.markAndSwap != 790 {
		t.Errorf("invalid wallTime.markAndSwap: got %v", i.wallTime.markAndSwap)
	}
	if i.wallTime.markTermination != 230 {
		t.Errorf("invalid wallTime.markTermination: got %v", i.wallTime.markTermination)
	}

	// CPU timings
	if i.cpuTime.sweepTermination != 130 {
		t.Errorf("invalid wallTime.sweepTermination: got %v", i.cpuTime.sweepTermination)
	}
	if i.cpuTime.markAndSwap != 4450 {
		t.Errorf("invalid wallTime.markAndSwap: got %v", i.cpuTime.markAndSwap)
	}
	if i.cpuTime.markTermination != 1800 {
		t.Errorf("invalid wallTime.markTermination: got %v", i.cpuTime.markTermination)
	}

	// Heap sizes
	if i.size.goal != 12 {
		t.Errorf("invalid size.goal: got %v", i.size.goal)
	}
	if i.size.live != 2 {
		t.Errorf("invalid size.goal: got %v", i.size.live)
	}
}

func TestParserGo15(t *testing.T) {
	line := "gc 24 @0.255s 11%: 0.073+0.19+0.46+0.77+0.33 ms clock, 0.58+0.19+0+0.24/0.92/0.88+2.6 ms cpu, 4->4->1 MB, 4 MB goal, 8 P"

	_, err := extractInfo(line)
	if err.Error() != "bad wall-clock timings" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParserGo14(t *testing.T) {
	line := "gc24(1): 0+0+800+0 us, 0 -> 1 MB, 10942 (60118-49176) objects, 281 goroutines, 216/0/0 sweeps, 0(0) handoff, 0(0) steal, 0/0/0 yields"

	_, err := extractInfo(line)
	if err.Error() != "bad input gctrace line" {
		t.Fatalf("unexpected error: %v", err)
	}
}
