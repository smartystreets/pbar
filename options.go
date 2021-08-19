package pbar

import (
	"io"
	"time"
)

// Option is a func type received by PBar.
// Each one allows configuration of the PBar.
type Option func(*PBar)

func RefreshIntervalMilliseconds(interval int) Option {
	return func(c *PBar) { c.refreshIntervalMilliseconds = time.Duration(interval) }
}

func BarLength(length int) Option {
	return func(c *PBar) { c.barLength = length }
}

func BarLeft(left rune) Option {
	return func(c *PBar) { c.barLeft = left }
}

func BarRight(right rune) Option {
	return func(c *PBar) { c.barRight = right }
}

func BarUncompleted(uncompleted rune) Option {
	return func(c *PBar) { c.barUncompleted = uncompleted }
}

func BarCompleted(completed rune) Option {
	return func(c *PBar) { c.barCompleted = completed }
}

func BarLabel(label string) Option {
	return func(c *PBar) { c.barLabel = label }
}

func OutputWriter(writer io.Writer) Option {
	return func(c *PBar) { c.output = writer }
}
