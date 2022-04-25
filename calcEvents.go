package main

import "math"

func calcEvents(lattice []Position, index int) {
	pos := lattice[index]
	esite := 0.0
	const EnergyLim = 10.0
	eventIndex := 0

	SCount, WCount := calcNeighbour(lattice, index)
	SNeighbours := float64(SCount)
	WNeighbours := float64(WCount)

	if pos.status != Hol {
		// Hole adsorption, atom desorption process
		// calc esite of current position
		if pos.status == Tug {
			esite = W_Ads - SNeighbours*WS_Bond - WNeighbours*WW_Bond
		} else if pos.status == Sul {
			if pos.Z == 0.0 {
				esite = S_Ads - WNeighbours*WS_Bond - SNeighbours*SS_Bond
			} else {
				esite = -WNeighbours*WS_Bond - SNeighbours*SS_Bond
			}
		}
		if -esite < EnergyLim {
			pos.events[eventIndex] = HolAds
			pos.rates[eventIndex] = v0 * math.Exp(esite*ev2j/kt)
			eventIndex++
		}
	} else {
		// Hole Desorption or Diffusion,  Atom Adsorption or Diffustion Process

		// Hole Diffusion process
		var nIndex int
		var deltaE, energy float64

		// Loop over diffusion path of 6 lateral neighbour
		for y := -1; y < 2; y++ {
			for x := -1; x < 2; x++ {
				if y != x {
					nX := int(pos.X) + x
					nY := int(pos.Y) + y
					nIndex = index + (y*lsize+x)*3
					nPos := lattice[nIndex]
					if nX < 1 || nY < 1 || nX > size || nY > size || nPos.status == Hol {
						// out of boundary or no neighbours, not allowed to happen
						continue
					}
					nSCount, nWCount := calcNeighbour(lattice, nIndex)
					nSNeighbours := float64(nSCount)
					nWNeighbours := float64(nWCount)
					if nPos.status == Sul {
						deltaE = (nWNeighbours-WNeighbours)*WS_Bond + (nSNeighbours-SNeighbours)*SS_Bond
						energy = math.Max(S_EDif, S_EDif+deltaE)
					} else {
						deltaE = (nSNeighbours-SNeighbours)*WS_Bond + (nWNeighbours-WNeighbours)*WW_Bond
						energy = math.Max(W_EDif, W_EDif+deltaE)
					}
					if energy < EnergyLim {
						var e Event
						switch c := [2]int{x, y}; c {
						case [2]int{0, -1}:
							e = HolDiffBL
						case [2]int{1, -1}:
							e = HolDiffBR
						case [2]int{-1, 0}:
							e = HolDiffL
						case [2]int{1, 0}:
							e = HolDiffR
						case [2]int{-1, 1}:
							e = HolDiffFL
						case [2]int{0, 1}:
							e = HolDiffFR
						}
						pos.events[eventIndex] = e
						pos.rates[eventIndex] = v0 * math.Exp(-energy*ev2j/kt)
						eventIndex++
					}
				}
			}
		}

		// Diffusion path of hole in vertical direction
		if pos.Z == 0.0 && lattice[index+2].status != Hol {
			nSCount, nWCount := calcNeighbour(lattice, index+2)
			nSNeighbours := float64(nSCount)
			nWNeighbours := float64(nWCount)
			deltaE = (nWNeighbours-WNeighbours)*WS_Bond + (nSNeighbours-SNeighbours)*SS_Bond + S_Ads
			energy = math.Max(S_EDif, S_EDif+deltaE)
			if energy < EnergyLim {
				pos.events[eventIndex] = HolDiffU
				pos.rates[eventIndex] = v0 * math.Exp(-energy*ev2j/kt)
				eventIndex++
			}
		}
		if pos.Z == 0.28 && lattice[index-2].status != Hol {
			nSCount, nWCount := calcNeighbour(lattice, index-2)
			nSNeighbours := float64(nSCount)
			nWNeighbours := float64(nWCount)
			deltaE = (nWNeighbours-WNeighbours)*WS_Bond + (nSNeighbours-SNeighbours)*SS_Bond - S_Ads
			energy = math.Max(S_EDif, S_EDif+deltaE)
			if energy < EnergyLim {
				pos.events[eventIndex] = HolDiffD
				pos.rates[eventIndex] = v0 * math.Exp(-energy*ev2j/kt)
				eventIndex++
			}
		}

		// Hole Desorption process
		if pos.Z == 0.14 {
			if SCount > 0 || WCount > 0 {
				energy := float64(SCount)*WS_Bond + float64(WCount)*WW_Bond
				pos.events[eventIndex] = HolDes
				pos.rates[eventIndex] = FluxRate * energy
				eventIndex++
			}
		}
		if pos.Z == 0.28 {
			if WCount > 0 || SCount > 0 {
				energy := float64(SCount)*SS_Bond + float64(WCount)*WS_Bond
				pos.events[eventIndex] = HolDes
				pos.rates[eventIndex] = FluxRate * CMRatio * energy
				eventIndex++
			}
		}
		if pos.Z == 0.0 {
			if WCount > 0 || SCount > 0 {
				energy := float64(SCount)*SS_Bond + float64(WCount)*WS_Bond
				pos.events[eventIndex] = HolDes
				pos.rates[eventIndex] = FluxRate * CMRatio * energy
				eventIndex++
			}
		}
	}

	if pos.sum != 0.0 {
		panic("Didn't clear previous result")
	}
	for _, v := range pos.rates {
		pos.sum += v
	}
	lattice[index] = pos
}
