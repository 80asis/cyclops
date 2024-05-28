package Task

type Task struct {
	Uuid []byte
}

func (t *Tasks) Queue(task Task) {
	// queues the task stack
}

func (t *Tasks) Dequeue() (task Task) {
	// dequeues the tasks stack

	return Task{
		Uuid: nil,
	}
}

type Tasks struct {
	TaskQueue []Task
}
