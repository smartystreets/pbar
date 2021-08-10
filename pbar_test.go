package pbar

import (
	"testing"

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
		BarLeft('<'), BarRight('>'), BarUncompleted('-'), BarCompleted('+'))

	this.So(progressBar.refreshIntervalMilliseconds, should.Equal, 750)
	this.So(progressBar.barLength, should.Equal, 25)
	this.So(progressBar.barLeft, should.Equal, '<')
	this.So(progressBar.barRight, should.Equal, '>')
	this.So(progressBar.barUncompleted, should.Equal, '-')
	this.So(progressBar.barCompleted, should.Equal, '+')
}

//func (this *PBarFixture) TestStart() {
//	progressBar := NewPBar(100, RefreshIntervalMilliseconds(1), BarLength(5))
//	progressBar.Start()
//
//	outBuf := bytes.NewBuffer(make([]byte, 0, 20))
//	progressBar.output = outBuf
//	progressBar.Update(100)
//	time.Sleep(5)
//	this.So(bytes.Runes(outBuf.Bytes()), should.Resemble, []rune("\x0D[=====] (100/100) 100% "))
//	progressBar.Finish()
//}
