package pbar

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
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

	this.So(progressBar.refreshInterval, should.Equal, 750*time.Millisecond)
	this.So(progressBar.barLength, should.Equal, 25)
	this.So(progressBar.barLeft, should.Equal, '<')
	this.So(progressBar.barRight, should.Equal, '>')
	this.So(progressBar.barUncompleted, should.Equal, '-')
	this.So(progressBar.barCompleted, should.Equal, '+')
}

func (this *PBarFixture) TestStart() {
	outBuf := bytes.NewBuffer(make([]byte, 0, 20))
	// Setting the output writer allows us to send all output to a buffer which allows us to test
	// that progressBar is creating output as expected.
	progressBar := NewPBar(1000, OutputWriter(outBuf),
		RefreshIntervalMilliseconds(250), BarLength(5))
	progressBar.Start()

	//safe read require to avoid race condition with refresh
	safeRead := func() []rune {
		progressBar.mutex.Lock()
		defer progressBar.mutex.Unlock()
		return bytes.Runes(outBuf.Bytes())
	}

	progressBar.Update(500)
	time.Sleep(time.Millisecond * 100)
	this.So(safeRead(), should.Resemble, []rune("\x0D[     ] (0/1,000) 0% "))

	time.Sleep(time.Millisecond * 300)
	this.So(safeRead(), should.Resemble,
		[]rune("\x0D[     ] (0/1,000) 0% \x0D[==   ] (500/1,000) 50% "))

	progressBar.Finish()

	//time.Sleep(time.Millisecond * 300)
	this.So(safeRead(), should.Resemble,
		[]rune("\x0D[     ] (0/1,000) 0% \x0D[==   ] (500/1,000) 50% \x0D[=====] (1,000/1,000) 100% \x0D[=====] (1,000/1,000) 100% "))
}

func (this *PBarFixture) TestCountFileLines() {
	tempFile := "/tmp/tempfille"
	fileContents := []byte("Line1\nLine2\nLine3\n")
	err := os.WriteFile(tempFile, fileContents, 0777)
	this.So(err, should.BeNil)
	defer func() { _ = os.Remove(tempFile) }()

	lineCount, err := CountFileLines(tempFile)
	this.So(err, should.BeNil)
	this.So(lineCount, should.Equal, 3)
}

func (this *PBarFixture) TestNoTerminal() {
	progressBar := NewPBar(1000,
		RefreshIntervalMilliseconds(250), BarLength(5))
	// Setting tty to a value other than the default (/dev/tty) will cause the tty open to fail,
	// simulating what occurs when a process using pbar is run in the background with no tty available like
	// when run with cron. In this case, all output is sent to io.Discard.
	//
	// Because no output is generated, this test demonstrates that a failure to
	// open the tty does not affect the calling program (like this test) from operating normally.
	progressBar.tty = "FALSE"
	progressBar.Start()
	progressBar.Update(250)
	time.Sleep(time.Millisecond * 250)
	progressBar.Update(500)
	time.Sleep(time.Millisecond * 250)
	progressBar.Update(750)
	progressBar.Finish()
	this.So(progressBar.currentCount, should.Equal, progressBar.TargetCount)
}
