package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/pbar"
)

func main() {
	fmt.Println()

	// get a new progress bar with all possible options
	progress := pbar.NewPBar(8000,
		pbar.BarLabel("File 1: "), pbar.RefreshIntervalMilliseconds(750), pbar.BarLength(25),
		pbar.BarLeft('⁅'), pbar.BarRight('⁆'), pbar.BarUncompleted('▭'), pbar.BarCompleted('▬'))
	// start the progress bar thread which updates the bar at the refresh interval
	progress.Start()
	fmt.Println()

	// get a second progress bar with all possible options
	progress2 := pbar.NewPBar(5000,
		pbar.BarLabel("File 2: "), pbar.RefreshIntervalMilliseconds(750), pbar.BarLength(25),
		pbar.BarLeft('⁅'), pbar.BarRight('⁆'), pbar.BarUncompleted('▭'), pbar.BarCompleted('▬'))
	progress2.Start()

	// simulate doing some stuff
	for i := 0; i <= 8000; i++ {
		if uint64(i) <= progress.TargetCount {
			progress.Update(uint64(i))       // update the counter in the progress bar
			time.Sleep(time.Millisecond / 2) // make it look like we are doing something important
		}

		if uint64(i) <= progress2.TargetCount {
			progress2.Update(uint64(i))
			time.Sleep(time.Millisecond / 2)
		}
	}

	progress.Finish() // update progress bar for the final time and terminate the thread
	progress2.Finish()
}
