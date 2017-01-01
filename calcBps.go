package main

const normalCo2 = 1800.0
const maxCo2    = 2500.0

const maxBps = 1.0
const minBps = 30.0

func calcBps(co2 float64) float64 {
	bpm := 0.0

	if co2 > maxCo2 {
		bpm = maxBps
	} else if co2 >= normalCo2  {
		bpm = maxBps + ((minBps-maxBps) / (maxCo2-normalCo2)) * (maxCo2 - co2)
 	}

	return bpm
}
