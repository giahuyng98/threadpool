package threadpool

type Task interface {
	Process() error
}
