package main

import (
	"time"

	"bitbucket.org/smartybryan/pbar"
)

func main() {
	// get a new progress bar with all possible options
	progress := pbar.NewPBar(5000,
		pbar.RefreshIntervalMilliseconds(750), pbar.BarLength(25),
		pbar.BarLeft('⁅'), pbar.BarRight('⁆'), pbar.BarUncompleted('▭'),pbar.BarCompleted('▬'))
	// start the progress bar thread which updates the bar at the refresh interval
	progress.Start()

	// simulate doing some stuff
	for i := 0; i <= 5000; i++ {
		progress.Update(uint64(i))   // update the counter in the progress bar
		time.Sleep(time.Millisecond) // make it look like we are doing something important
	}

	progress.Finish() // update progress bar for the final time and terminate the thread
}
