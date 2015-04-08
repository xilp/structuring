package slave

type Task interface {
	Url () string
	Run(chan Task) error
	Weight() int
}

type Site interface {
	NewTask(info TaskInfo) Task
}
