package pbar

import (
	"io"
	"time"
)

// Option is a func type received by PBar.
// Each one allows configuration of the PBar.
type Option func(*PBarSupport)

func RefreshIntervalMilliseconds(interval int) Option {
	return func(c *PBarSupport) { c.refreshInterval = time.Duration(interval) * time.Millisecond }
}

func BarLength(length int) Option {
	return func(c *PBarSupport) { c.barLength = length }
}

func BarLeft(left rune) Option {
	return func(c *PBarSupport) { c.barLeft = left }
}

func BarRight(right rune) Option {
	return func(c *PBarSupport) { c.barRight = right }
}

func BarUncompleted(uncompleted rune) Option {
	return func(c *PBarSupport) { c.barUncompleted = uncompleted }
}

func BarCompleted(completed rune) Option {
	return func(c *PBarSupport) { c.barCompleted = completed }
}

func BarLabel(label string) Option {
	return func(c *PBarSupport) { c.barLabel = label }
}

func OutputWriter(writer io.Writer) Option {
	return func(c *PBarSupport) {
		c.testing = true
		c.output = writer
	}
}
