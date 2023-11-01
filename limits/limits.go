package limits

var DefaultLimits = Limits{
	CPUPercentage: 1.0,
	MemoryKB:      -1,
}

type Limits struct {
	CPUPercentage float64
	MemoryKB      int64
}

type LimitOption func(*Limits)

func WithCPUPercentage(p float64) LimitOption {
	return func(l *Limits) {
		l.CPUPercentage = p
	}
}

func WithMemoryKB(m int64) LimitOption {
	return func(l *Limits) {
		l.MemoryKB = m
	}
}
