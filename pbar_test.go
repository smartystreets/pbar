package pbar

import (
	"bytes"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPBarFixture(t *testing.T) {
	gunit.Run(new(PBarFixture), t)
}

type PBarFixture struct {
	*gunit.Fixture
}

func (this *PBarFixture) TestOptions() {
	progressBar := NewPBar(100,
		RefreshIntervalMilliseconds(750), BarLength(25),
		BarLeft('<'), BarRight('>'), BarCompleted('+'))

	this.So(progressBar.refreshIntervalMilliseconds, should.Equal, 750)
	this.So(progressBar.barLength, should.Equal, 25)
	this.So(progressBar.barLeft, should.Equal, '<')
	this.So(progressBar.barRight, should.Equal, '>')
	this.So(progressBar.barCompleted, should.Equal, '+')
}

func (this *PBarFixture) TestStart() {
	progressBar := NewPBar(100, RefreshIntervalMilliseconds(1), BarLength(5))
	progressBar.Start()

	outBuf := bytes.NewBuffer(make([]byte, 0, 20))
	progressBar.output = outBuf
	progressBar.Update(100) // signal the completion (current == targetCount)
	time.Sleep(time.Millisecond) // wait for update thread to terminate
	this.So(outBuf.Bytes(), should.Resemble, []byte("\x0D[-----] (100/100) 100% "))
}

func (this *PBarFixture) TestUpdate() {
	progressBar := NewPBar(1000, BarLength(10))
	progressBar.Update(100)
	progressBar.updateBar()
	this.So(progressBar.currentCount, should.Equal, 100)
	this.So(progressBar.barVisual, should.Resemble, []byte("[-         ]"))
	this.So(progressBar.barPercent, should.Equal, "(100/1,000) 10%")
}
