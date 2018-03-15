package mail

type Logger interface {
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
}
