package pbar

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	BarLengthDefault       = 50
	RefreshIntervalDefault = 500
	BarLeftDefault         = '['
	BarRightDefault        = ']'
	BarUnCompletedDefault  = ' '
	BarCompletedDefault    = '='
)

type PBar struct {
	barVisual    []rune
	barPercent   string
	currentCount uint64
	targetCount  uint64
	output       io.Writer

	refreshIntervalMilliseconds time.Duration
	barLength                   int
	barLeft                     rune
	barRight                    rune
	barUncompleted              rune
	barCompleted                rune
}

func NewPBar(targetCount uint64, options ...Option) *PBar {
	return new(PBar).configure(targetCount, options)
}

func (this *PBar) configure(targetCount uint64, options []Option) *PBar {
	this.targetCount = targetCount
	this.output = os.Stdout

	this.barLength = BarLengthDefault
	this.refreshIntervalMilliseconds = RefreshIntervalDefault
	this.barLeft = BarLeftDefault
	this.barRight = BarRightDefault
	this.barUncompleted = BarUnCompletedDefault
	this.barCompleted = BarCompletedDefault

	for _, configure := range options {
		configure(this)
	}

	return this
}

func (this *PBar) Start() {
	var waiter sync.WaitGroup
	waiter.Add(1)
	go this.start(&waiter)
	waiter.Wait()
}

func (this *PBar) start(waiter *sync.WaitGroup) {
	this.initializeBar()
	waiter.Done()

	for {
		if this.currentCount == this.targetCount {
			this.updateBar()
			this.repaint()
			break
		}
		this.updateBar()
		this.repaint()
		time.Sleep(time.Millisecond * this.refreshIntervalMilliseconds)
	}
}

func (this *PBar) Finish() {
	this.currentCount = this.targetCount
	this.updateBar()
	this.repaint()
}

func (this *PBar) Update(current uint64) {
	this.currentCount = current
}

func (this *PBar) initializeBar() {
	this.barVisual = make([]rune, this.barLength+2) // plus beginning and end markers
	this.barVisual[0] = this.barLeft
	this.barVisual[this.barLength+1] = this.barRight
	this.updateBar()
}

func (this *PBar) updateBar() {
	percentCompleted := float32(this.currentCount) / float32(this.targetCount)
	completed := int(percentCompleted * float32(this.barLength))

	for i := 1; i <= this.barLength; i++ {
		if i <= completed {
			this.barVisual[i] = this.barCompleted
		} else {
			this.barVisual[i] = this.barUncompleted
		}
	}

	this.barPercent = fmt.Sprintf("(%s/%s) %d%%",
		comma(this.currentCount), comma(this.targetCount), int(percentCompleted*100.0))
}

func (this *PBar) repaint() {
	// go to beginning of the line and print data
	_, _ = fmt.Fprintf(this.output, "%c%s %s%c", 13, string(this.barVisual), this.barPercent, 32)
}

func comma(n uint64) string {
	in := strconv.FormatUint(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
