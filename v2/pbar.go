package pbar

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/term"
)

const (
	TTY = "/dev/tty" // Microsoft Windows is not supported

	BarLengthDefault       = 50
	RefreshIntervalDefault = 500 * time.Millisecond
	BarLeftDefault         = '['
	BarRightDefault        = ']'
	BarUnCompletedDefault  = ' '
	BarCompletedDefault    = '='
)

type PBar[T integer] struct {
	PBarSupport
	mutex        sync.Mutex
	currentCount T
	TargetCount  T
}

type PBarSupport struct {
	barVisual      []rune
	barPercent     string
	terminal       *term.Term
	cursorPosition CursorPosition
	tty            string

	refreshInterval                                 time.Duration
	barLength                                       int
	barLeft, barRight, barUncompleted, barCompleted rune
	barLabel                                        string
	testing                                         bool
	output                                          io.Writer
}

func DefaultPBarSupport() PBarSupport {
	return PBarSupport{
		barLength:       BarLengthDefault,
		refreshInterval: RefreshIntervalDefault,
		barLeft:         BarLeftDefault,
		barRight:        BarRightDefault,
		barUncompleted:  BarUnCompletedDefault,
		barCompleted:    BarCompletedDefault,
		tty:             TTY,
		output:          os.Stdout,
	}
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type CursorPosition struct {
	row int8
	col int8
}

func NewPBar[T integer](targetCount T, options ...Option) *PBar[T] {
	return new(PBar[T]).configure(targetCount, options)
}

// [locks mutex]
func (this *PBar[T]) configure(targetCount T, options []Option) *PBar[T] {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.TargetCount = targetCount
	this.PBarSupport = DefaultPBarSupport()

	for _, configure := range options {
		configure(&this.PBarSupport)
	}

	return this
}

func (this *PBar[T]) Start() {
	var waiter sync.WaitGroup
	waiter.Add(1)
	go this.start(&waiter)
	waiter.Wait()
}

// [locks mutex]
func (this *PBar[T]) start(waiter *sync.WaitGroup) {
	this.saveCursorPosition()
	this.initializeBar()
	waiter.Done()

	for {
		this.updateBar()
		this.repaint()
		this.mutex.Lock()
		done := this.currentCount == this.TargetCount
		this.mutex.Unlock()

		if done {
			break
		}
		time.Sleep(this.refreshInterval)
	}
}

// [locks mutex]
func (this *PBar[T]) Finish() {
	this.mutex.Lock()
	this.currentCount = this.TargetCount
	this.mutex.Unlock()
	time.Sleep(this.refreshInterval)
	this.updateBar()
	this.repaint()
}

// [locks mutex]
func (this *PBar[T]) Update(current T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.currentCount = current
}

// [locks mutex]
func (this *PBar[T]) updateBar() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	percentCompleted := float32(this.currentCount) / float32(this.TargetCount)
	completed := int(percentCompleted * float32(this.barLength))

	for i := 1; i <= this.barLength; i++ {
		if i <= completed {
			this.barVisual[i] = this.barCompleted
		} else {
			this.barVisual[i] = this.barUncompleted
		}
	}

	this.barPercent = fmt.Sprintf("(%s/%s) %d%%",
		comma(this.currentCount), comma(this.TargetCount), int(percentCompleted*100.0))
}

// [locks mutex]
func (this *PBar[T]) repaint() {
	this.restoreCursorPosition()
	this.mutex.Lock()
	// go to beginning of the line and print data
	_, _ = fmt.Fprintf(this.output, "%c%s%s %s%c", 13, this.barLabel, string(this.barVisual), this.barPercent, 32)
	this.mutex.Unlock()
}

// [locks mutex]
func (this *PBar[T]) openTty() (err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.terminal, err = term.Open(this.tty)
	if err != nil {
		this.testing = true // prevent attempts to save and restore cursor position
		this.output = io.Discard
		return
	}
	_ = term.RawMode(this.terminal)

	return
}

// [locks mutex]
func (this *PBar[T]) closeTty() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	_ = this.terminal.Restore()
}

// [locks mutex]
func (this *PBar[T]) saveCursorPosition() {
	if this.testing {
		return
	}

	if this.openTty() != nil {
		return
	}
	defer this.closeTty()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	out := make([]byte, 6)
	_, _ = this.terminal.Write([]byte{13, 27, '[', '6', 'n'})
	_, _ = this.terminal.Read(out)
	split := strings.Split(string(out[2:]), ";")
	if len(split) > 1 {
		this.cursorPosition.row = atoi8(split[0])
		this.cursorPosition.col = atoi8(split[1])
	}
}

// [locks mutex]
func (this *PBar[T]) restoreCursorPosition() {
	if this.testing {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.cursorPosition.row == 0 && this.cursorPosition.col == 0 {
		return
	}
	fmt.Printf("%c%c%d;%dH", 27, '[', this.cursorPosition.row, this.cursorPosition.col)
}

// CountFileLines count newline characters in a file
func CountFileLines(path string) (count int, err error) {
	const lineBreak = '\n'

	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// [locks mutex]
func (this *PBar[T]) initializeBar() {
	this.mutex.Lock()
	this.barVisual = make([]rune, this.barLength+2) // plus beginning and end markers
	this.barVisual[0] = this.barLeft
	this.barVisual[this.barLength+1] = this.barRight
	this.mutex.Unlock()

	this.updateBar()
}

func atoi8(val string) int8 {
	strVal, _ := strconv.Atoi(val)
	return int8(strVal)
}

func comma[T integer](n T) string {
	in := fmt.Sprintf("%d", n)
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
