package main

import (
	"fmt"
	"math/rand"
	"os"
)

func recordEvent(lattice []Position, index int, subIndex int, recorder *EventRecorder) {
	switch pos := lattice[index]; pos.events[subIndex] {
	case HolDes:
		if pos.Z == 0.14 {
			recorder.WAds++
		} else {
			recorder.SAds++
		}
	case HolAds:
		if pos.Z == 0.14 {
			recorder.WDes++
		} else {
			recorder.SDes++
		}
	case HolDiffFR:
		if pos.Z == 0.14 {
			recorder.WDiffBL++
		} else {
			recorder.SDiffBL++
		}
	case HolDiffR:
		if pos.Z == 0.14 {
			recorder.WDiffL++
		} else {
			recorder.SDiffL++
		}
	case HolDiffBR:
		if pos.Z == 0.14 {
			recorder.WDiffFL++
		} else {
			recorder.SDiffFL++
		}
	case HolDiffBL:
		if pos.Z == 0.14 {
			recorder.WDiffFR++
		} else {
			recorder.SDiffFR++
		}
	case HolDiffL:
		if pos.Z == 0.14 {
			recorder.WDiffR++
		} else {
			recorder.SDiffR++
		}
	case HolDiffFL:
		if pos.Z == 0.14 {
			recorder.WDiffBR++
		} else {
			recorder.SDiffBR++
		}
	case HolDiffD:
		if pos.Z == 0.14 {
			recorder.WDiffU++
		} else {
			recorder.SDiffU++
		}
	case HolDiffU:
		if pos.Z == 0.14 {
			recorder.WDiffD++
		} else {
			recorder.WDiffD++
		}
	}
}

func chooseEvent(lattice []Position, sumRates float64) (int, int) {
	index := 0
	rand1 := rand.Float64() * sumRates
	posRates := 0.0
	for y := 1; y < size+1; y++ {
		for x := 1; x < size+1; x++ {
			for k := 0; k < 3; k++ {
				index = (y*lsize+x)*3 + k
				posRates += lattice[index].sum
				if posRates > rand1 {
					// find position to execute event
					pos := lattice[index]
					rand2 := rand.Float64() * pos.sum
					eventRates := 0.0
					for subIndex, subRate := range pos.rates {
						eventRates += subRate
						if eventRates > rand2 {
							// find event in that position
							return index, subIndex
						}
					}
				}
			}
		}
	}
	fmt.Println(index, posRates, sumRates)
	panic("No Event found")
}

func executeEvent(lattice []Position, index, subIndex int) {
	if lattice[index].events[subIndex] == HolAds {
		updateStatus(lattice, index, Hol)
	} else {
		if lattice[index].Z == 0.14 {
			updateStatus(lattice, index, Tug)
		} else {
			updateStatus(lattice, index, Sul)
		}
		switch lattice[index].events[subIndex] {
		case HolDiffFR:
			updateStatus(lattice, index+3, Hol)
		case HolDiffR:
			updateStatus(lattice, index+lsize*3, Hol)
		case HolDiffBR:
			updateStatus(lattice, index+lsize*3-3, Hol)
		case HolDiffBL:
			updateStatus(lattice, index-3, Hol)
		case HolDiffL:
			updateStatus(lattice, index-lsize*3, Hol)
		case HolDiffFL:
			updateStatus(lattice, index-lsize*3+3, Hol)
		case HolDiffD:
			updateStatus(lattice, index-2, Hol)
		case HolDiffU:
			updateStatus(lattice, index+2, Hol)
		}
	}
}

func clearEvents(lattice []Position) {
	for index := range lattice {
		lattice[index].events = [8]Event{}
		lattice[index].rates = [8]float64{}
		lattice[index].sum = 0.0
	}
}

func updateStatus(lattice []Position, index int, target Atom) {
	if pos := lattice[index]; pos.X < 1 || pos.X > size || pos.Y < 1 || pos.Y > size {
		return
	}
	indexA, indexB := checkBorder(&lattice[index])
	lattice[index].status = target
	lattice[index+indexA].status = target
	lattice[index+indexB].status = target
}

// return additional index at border that need to be updated
func checkBorder(pos *Position) (int, int) {
	if pos.Y == 1.0 {
		if pos.X == 1.0 {
			return size * lsize * 3, size * 3
		} else if pos.X == size+1 {
			return size * lsize * 3, -size * 3
		} else {
			return size * lsize * 3, 0
		}
	}
	if pos.Y == size+1 {
		if pos.X == 1.0 {
			return -size * lsize * 3, size * 3
		} else if pos.X == size+1 {
			return -size * lsize * 3, -size * 3
		} else {
			return -size * lsize * 3, 0
		}
	}
	if pos.X == 1.0 {
		return size * 3, 0
	}
	if pos.X == size+1 {
		return -size * 3, 0
	}
	return 0, 0
}

// Lattice Structure
//         (-1, 1, 0) ---- (0, 1, 0)
//         /        \      /       \
//        /          \    /         \
//       /            \  /           \
//  (-1, 0, 0) ---- (0, 0, 0) --- (1, 0, 0)
//       \            /  \           /
//        \          /    \         /
//         \        /      \       /
//         (-1, 1, 0) ---- (0, 1, 0)
func calcNeighbour(lattice []Position, index int) (int, int) {
	SNeighbour := 0
	WNeighbour := 0
	var nIndex int

	for y := -1; y < 2; y++ {
		for x := -1; x < 2; x++ {
			if y != x {
				nIndex = index + (y*lsize+x)*3
				if lattice[nIndex].status == Sul {
					SNeighbour++
				} else if lattice[nIndex].status == Tug {
					WNeighbour++
				}
			}
		}
	}

	if lattice[index].Z == 0.14 {
		if lattice[index+1].status == Sul {
			SNeighbour++
		}
		if lattice[index-1].status == Sul {
			SNeighbour++
		}
		if lattice[index+3+1].status == Sul {
			SNeighbour++
		}
		if lattice[index+3-1].status == Sul {
			SNeighbour++
		}
		if lattice[index+3*lsize+1].status == Sul {
			SNeighbour++
		}
		if lattice[index+3*lsize-1].status == Sul {
			SNeighbour++
		}
	}

	if lattice[index].Z == 0.0 {
		if lattice[index+1].status == Tug {
			WNeighbour++
		}
		if lattice[index-3+1].status == Tug {
			WNeighbour++
		}
		if lattice[index-3*lsize+1].status == Tug {
			WNeighbour++
		}
	} else {
		if lattice[index-1].status == Tug {
			WNeighbour++
		}
		if lattice[index-3-1].status == Tug {
			WNeighbour++
		}
		if lattice[index-3*lsize-1].status == Tug {
			WNeighbour++
		}
	}

	return SNeighbour, WNeighbour
}

func initLattice(size int) []Position {
	lattice := make([]Position, lsize*lsize*3)

	// initialize lattice
	index := 0
	for y := 0; y < lsize; y++ {
		for x := 0; x < lsize; x++ {
			index = (y*lsize + x) * 3
			lattice[index] = Position{float32(x), float32(y), 0.0, Hol, [8]Event{}, [8]float64{}, 0.0}
			lattice[index+1] = Position{float32(x) + 0.33, float32(y) + 0.33, 0.14, Hol, [8]Event{}, [8]float64{}, 0.0}
			lattice[index+2] = Position{float32(x), float32(y), 0.28, Hol, [8]Event{}, [8]float64{}, 0.0}
		}
	}

	// initialize point in lattice
	initSize := size / 5
	for y := 1 + (size-initSize)/2; y < 1+(size+initSize)/2; y++ {
		for x := 1 + (size-initSize)/2; x < 1+(size+initSize)/2; x++ {
			if x+y < 2+size {
				index = (y*lsize + x) * 3
				lattice[index] = Position{float32(x), float32(y), 0.0, Sul, [8]Event{}, [8]float64{}, 0.0}
				lattice[index+1] = Position{float32(x) + 0.33, float32(y) + 0.33, 0.14, Tug, [8]Event{}, [8]float64{}, 0.0}
				lattice[index+2] = Position{float32(x), float32(y), 0.28, Sul, [8]Event{}, [8]float64{}, 0.0}
			}
		}
	}

	return lattice
}

func writeToResult(lattice []Position) {
	f, _ := os.Create("result.py")
	defer f.Close()

	f.WriteString("result=[")
	var index int
	for y := 1; y < size+1; y++ {
		for x := 1; x < size+1; x++ {
			index = (y*lsize + x) * 3
			pos1 := lattice[index]
			pos2 := lattice[index+1]
			pos3 := lattice[index+2]
			if pos1.status == Sul || pos3.status == Sul {
				if pos1.status == Sul && pos3.status == Sul {
					// yellow for two sulphur atoms
					f.WriteString(fmt.Sprintf("[%v, %v, 'yellow'],\n", pos1.X+0.5*pos1.Y, 0.866*pos1.Y))
				} else {
					// orange for one sulphur atom
					f.WriteString(fmt.Sprintf("[%v, %v, 'orange'],\n", pos1.X+0.5*pos1.Y, 0.866*pos1.Y))
				}
			}
			if pos2.status == Tug {
				f.WriteString(fmt.Sprintf("[%v, %v, 'blue'],\n", pos2.X+0.5*pos2.Y, 0.866*pos2.Y))
			}
		}
	}
	f.WriteString("]\n")
	f.WriteString(fmt.Sprintf("xlim = %v\n", size+0.5*size))
	f.WriteString(fmt.Sprintf("ylim = %v\n", 0.866*size))
}
