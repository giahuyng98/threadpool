package threadpool

import (
	"fmt"
	"sync"
)

type Worker struct {
	id    int
	tasks chan Task
	wg    *sync.WaitGroup
}

func (w *Worker) Start() {
	defer w.wg.Done()
	for {
		task, ok := <-w.tasks
		if !ok {
			return
		}
		if task == nil {
			continue
		}
		task.Process()
	}
}

//
type SafeQueue struct {
	data []Task
	mu   sync.RWMutex
}

func (q *SafeQueue) Push(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, task)
}

func (q *SafeQueue) Peek() (Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil, fmt.Errorf("Queue is empty")
	}
	result := q.data[0]
	return result, nil
}

func (q *SafeQueue) Pop() (Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil, fmt.Errorf("Queue is empty")
	}
	result := q.data[0]
	q.data = q.data[1:]
	return result, nil
}

func (q *SafeQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data)
}

func (q *SafeQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data) == 0
}

type Pool struct {
	size        int
	tasks       chan Task
	workerTasks chan Task
	idleTasks   SafeQueue
	workers     []Worker
	wg          sync.WaitGroup
	quit        chan bool
}

func (p *Pool) AddTask(task Task) {
	//p.idleTasks.Push(task)
	p.tasks <- task
}

// TODO: quit it
func (p *Pool) Join() {
	<-p.quit
	close(p.workerTasks)
	p.wg.Wait()
}

func New(poolSize, queueCap int) *Pool {
	if poolSize <= 0 {
		poolSize = 1
	}
	pool := &Pool{
		size:        poolSize,
		workerTasks: make(chan Task),
		tasks:       make(chan Task),
		idleTasks: SafeQueue{
			data: make([]Task, 0, queueCap),
		},
		workers: make([]Worker, poolSize),
		quit:    make(chan bool),
	}
	pool.wg.Add(poolSize + 1)

	for i, _ := range pool.workers {
		pool.workers[i] = Worker{
			id:    i,
			tasks: pool.workerTasks,
			wg:    &pool.wg,
		}
		go pool.workers[i].Start()
	}
	go pool.run()
	return pool
}

func (p *Pool) run() {
	defer p.wg.Done()

	for {
		task, err := p.idleTasks.Peek()
		if err != nil {
			select {
			case p.workerTasks <- <-p.tasks:
			case t := <-p.tasks:
				p.idleTasks.Push(t)
			}
		} else {
			select {
			case p.workerTasks <- task:
				p.idleTasks.Pop()
			case t := <-p.tasks:
				p.idleTasks.Push(t)
			}
		}
	}
}
