package logger

type StubLogger struct {
}

func (l *StubLogger) Debug(_ string, _ ...any) {
}

func (l *StubLogger) Info(_ string, _ ...any) {
}

func (l *StubLogger) Warn(_ string, _ ...any) {
}

func (l *StubLogger) Error(_ string, _ ...any) {
}
