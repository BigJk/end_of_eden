package audio

type emptySound struct{}

func (e emptySound) Len() int {
	return sampleRate
}

func (e emptySound) Position() int {
	return 0
}

func (e emptySound) Seek(p int) error {
	return nil
}

func (e emptySound) Stream(samples [][2]float64) (n int, ok bool) {
	return len(samples), true
}

func (e emptySound) Err() error {
	return nil
}
