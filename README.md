# SMARTY DISCLAIMER: Subject to the terms of the associated license agreement, this software is freely available for your use. This software is FREE, AS IN PUPPIES, and is a gift. Enjoy your new responsibility. This means that while we may consider enhancement requests, we may or may not choose to entertain requests at our sole and absolute discretion.

Progress Bar (PBar)
============================

The Progress Bar (PBar) Go library provides basic text progress bar functionality.

#### Features

* Configurable for progress bar length and characters used to paint the bar
* Multi-threaded to paint independently of the underlying looping process

## Import
```
import github.com/smartystreets/pbar
```

## Usage
```
progress := pbar.NewPBar(5000) // 5000 is the target count
progress.Start()

... some looping work to do
progress.Update(counter)
... end of loop

progress.Finish()
``` 

## Options
Specify any number of comma separated options as parameters to `NewPBar()`

For example:
```
progress := pbar.NewPBar(5000, pbar.BarLength(25), pbar.RefreshIntervalMilliseconds(750))
```

#### Progress Bar Label
Set the label that will be output to the left of the bar.
```
pbar.BarLabel("Loading index: ")
```

#### Progress Bar Length
Set the length of the progress bar, not counting the summary text. Default 50.
```
pbar.BarLength(25)
```

#### Progress Bar Refresh Interval
Set the refresh interval of the progress bar in milliseconds.  Default 500ms.
```
pbar.RefreshIntervalMilliseconds(750)
```

#### Progress Bar Graphic Characters
Set the graphic characters used in the progress bar.
```
pbar.BarLeft('⁅')
pbar.BarRight('⁆')
pbar.BarUncompleted('▭')
pbar.BarCompleted('▬')
```

## Example Code
See `cmd/main.go` for a fully functional sample.
