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
		BarLabel("Testing"), RefreshIntervalMilliseconds(750), BarLength(25),
		BarLeft('<'), BarRight('>'), BarUncompleted('-'), BarCompleted('+'))

	this.So(progressBar.refreshIntervalMilliseconds, should.Equal, 750)
	this.So(progressBar.barLength, should.Equal, 25)
	this.So(progressBar.barLeft, should.Equal, '<')
	this.So(progressBar.barRight, should.Equal, '>')
	this.So(progressBar.barUncompleted, should.Equal, '-')
	this.So(progressBar.barCompleted, should.Equal, '+')
}

func (this *PBarFixture) TestStart() {
	outBuf := bytes.NewBuffer(make([]byte, 0, 20))
	progressBar := NewPBar(1000, OutputWriter(outBuf),
		RefreshIntervalMilliseconds(250), BarLength(5))
	progressBar.Start()

	progressBar.Update(500)
	time.Sleep(time.Millisecond * 100)
	this.So(bytes.Runes(outBuf.Bytes()), should.Resemble, []rune("\x0D[     ] (0/1,000) 0% "))

	time.Sleep(time.Millisecond * 300)
	this.So(bytes.Runes(outBuf.Bytes()), should.Resemble,
		[]rune("\x0D[     ] (0/1,000) 0% \x0D[==   ] (500/1,000) 50% "))

	progressBar.Finish()

	//time.Sleep(time.Millisecond * 300)
	this.So(bytes.Runes(outBuf.Bytes()), should.Resemble,
		[]rune("\x0D[     ] (0/1,000) 0% \x0D[==   ] (500/1,000) 50% \x0D[=====] (1,000/1,000) 100% "))
}
