package main

import "fmt"

const (
	size  = 50
	lsize = size + 2 // for periodicity, makes lattice larger
	step  = 1e4
)

const (
	WW_Bond = 1.16
	WS_Bond = 2.14
	SS_Bond = 0.329
	W_Ads   = -2.0
	S_Ads   = -1.28
	SDif    = 0.03
	W_EDif  = 3.81
	S_EDif  = 1.02

	T    = 273.15 + 1000
	k    = 1.380649e-23
	ev2j = 1.6e-19
	kt   = k * T
	h    = 6.62607015e-34
	v0   = 2.0 * kt / h

	FluxRate = 1e3
	CMRatio  = 2.0
)

type Atom int8

const (
	Sul Atom = iota
	Tug
	Hol
)

func (s Atom) String() string {
	switch s {
	case Sul:
		return "S"
	case Tug:
		return "W"
	case Hol:
		return "Hole"
	}
	return "unknown"
}

type Event int8

const (
	// Use hole to simplify process
	HolDes Event = iota + 1 // 0 used to represent none, start from 0
	HolAds
	// Direction: front, right, back, left + up, down
	HolDiffFR
	HolDiffR
	HolDiffBR
	HolDiffBL
	HolDiffL
	HolDiffFL
	HolDiffD
	HolDiffU
)

func (e Event) String() string {
	switch e {
	case HolDes:
		return "HolDes"
	case HolAds:
		return "HolAds"
	case HolDiffFR:
		return "HolDiffFR"
	case HolDiffR:
		return "HolDiffR"
	case HolDiffBR:
		return "HolDiffBR"
	case HolDiffBL:
		return "HolDiffBL"
	case HolDiffL:
		return "HolDiffL"
	case HolDiffFL:
		return "HolDiffFL"
	case HolDiffD:
		return "HolDiffD"
	case HolDiffU:
		return "HolDiffU"
	}
	return "none"
}

type Position struct {
	X, Y, Z float32
	status  Atom
	events  [8]Event
	rates   [8]float64
	sum     float64
}

func (p Position) String() string {
	return fmt.Sprintf("%v, %v, %v, %v, %v, %v, %v\n", p.X, p.Y, p.Z, p.status, p.events, p.rates, p.sum)
}

type EventRecorder struct {
	WAds, WDes, WDiffFR, WDiffR, WDiffBR, WDiffBL, WDiffL, WDiffFL, WDiffD, WDiffU uint64
	SAds, SDes, SDiffFR, SDiffR, SDiffBR, SDiffBL, SDiffL, SDiffFL, SDiffD, SDiffU uint64
}

func (e EventRecorder) String() string {
	return fmt.Sprintf("WAds: %v\nWDes: %v\nWDiffFR: %v\nWDiffR: %v\nWDiffBR: %v\nWDiffBL: %v\nWDiffL: %v\nWDiffFL: %v\nWDiffD: %v\nWDiffU: %v\nSAds: %v\nSDes: %v\nSDiffFR: %v\nSDiffR: %v\nSDiffBR: %v\nSDiffBL: %v\nSDiffL: %v\nSDiffFL: %v\nSDiffD: %v\nSDiffU: %v", e.WAds, e.WDes, e.WDiffFR, e.WDiffR, e.WDiffBR, e.WDiffBL, e.WDiffL, e.WDiffFL, e.WDiffD, e.WDiffU, e.SAds, e.SDes, e.SDiffFR, e.SDiffR, e.SDiffBR, e.SDiffBL, e.SDiffL, e.SDiffFL, e.SDiffD, e.SDiffU)
}
