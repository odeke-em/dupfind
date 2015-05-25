package dupfinder

import (
	spinner "github.com/odeke-em/cli-spinner"
)

type playable struct {
	Play  func()
	Pause func()
	Reset func()
	Stop  func()
}

func noop() {
}

func noopPlayable() *playable {
	return &playable{
		Play:  noop,
		Pause: noop,
		Reset: noop,
		Stop:  noop,
	}
}

func NewPlayable(freq int64) *playable {
	spin := spinner.New(freq)

	play := func() {
		spin.Start()
	}

	return &playable{
		Play:  play,
		Stop:  spin.Stop,
		Reset: spin.Reset,
		Pause: spin.Stop,
	}
}
