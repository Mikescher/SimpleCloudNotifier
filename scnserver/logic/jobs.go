package logic

type Job interface {
	Start() error
	Stop()
	Running() bool
}
