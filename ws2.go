package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	start := time.Now()
	// rand.Seed(time.Now().UnixNano())

	lattice := initLattice(size)
	index := 0
	reactionTime := 0.0
	recorder := EventRecorder{}

	for t := 0; t < int(step); t++ {
		// tabulate all possible events
		sumRates := 0.0
		for y := 1; y < size+1; y++ {
			for x := 1; x < size+1; x++ {
				index = (y*lsize + x) * 3
				calcEvents(lattice, index)
				calcEvents(lattice, index+1)
				calcEvents(lattice, index+2)
				sumRates += lattice[index].sum
				sumRates += lattice[index+1].sum
				sumRates += lattice[index+2].sum
			}
		}

		// fmt.Println(sumRates)

		// choose an event
		if sumRates > 0.0 {
			index, subIndex := chooseEvent(lattice, sumRates)
			reactionTime += 1.0 / lattice[index].rates[subIndex]
			executeEvent(lattice, index, subIndex)
			recordEvent(lattice, index, subIndex, &recorder)
			clearEvents(lattice)
		}

		if t > 0 && t%500 == 0 {
			fmt.Println("KMC Simulation: ", t, "steps")
		}
	}

	fmt.Println(recorder)

	writeToResult(lattice)

	elapsed := time.Since(start)
	fmt.Printf("Runtime: %s", elapsed)

	exec.Command("python", "plot.py").Run()
}
