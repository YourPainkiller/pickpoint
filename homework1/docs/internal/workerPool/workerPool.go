package workerPool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Структура воркер пула. Некоторые поля полноценно не используются, но я подумал что они могут быть полезными в будущем
type Pool struct {
	workerCount         int32
	waitingTaskCount    int32
	succesfullTaskCount int32
	busyWorkerCount     int32
	freeWorkerCount     int32

	maxWorkers int
	maxTasks   int
	tasksWg    sync.WaitGroup
	workersWg  sync.WaitGroup //queue of workers to be deleted

	tasks chan func()
	quit  chan struct{}
}

// Создание нового пула
func NewPool(maxWorkers, maxTasks int) *Pool {
	pool := &Pool{
		maxWorkers: maxWorkers,
		maxTasks:   maxTasks,
	}

	if pool.maxWorkers < 1 {
		pool.maxWorkers = 1
	}
	if pool.maxTasks < 1 {
		pool.maxTasks = 1
	}

	pool.tasks = make(chan func(), maxTasks)
	pool.quit = make(chan struct{}, maxWorkers)

	return pool
}

func (p *Pool) GetTasksWg() *sync.WaitGroup {
	return &p.tasksWg
}

func (p *Pool) GetWorkersWg() *sync.WaitGroup {
	return &p.workersWg
}

func (p *Pool) GetCurrentWorkers() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *Pool) GetMaxWorkers() int {
	return p.maxWorkers
}

// Создание воркера
func (p *Pool) CreateWorker() {
	totalWorkers := atomic.LoadInt32(&p.workerCount)
	if totalWorkers >= int32(p.maxWorkers) {
		//fmt.Print("Can't create new worker, because maxWorkers reached\n")
		return
	}
	id := totalWorkers + 1
	atomic.AddInt32(&p.workerCount, 1)
	atomic.AddInt32(&p.freeWorkerCount, 1)
	fmt.Printf("starting worker%d\n", id)

	for {
		select {
		case task := <-p.tasks:
			//fmt.Printf("Worker%d: ", id)
			p.startTask()
			task()
			p.endTask()
			p.tasksWg.Done()
		case <-p.quit:
			p.workersWg.Done()
			atomic.AddInt32(&p.workerCount, -1)
			atomic.AddInt32(&p.freeWorkerCount, -1)
			fmt.Printf("worker%d deleted\n", id)
			return
		}
	}
}

// Добавление задачи в воркер пул
func (p *Pool) SubmitTask(task func()) {
	p.tasksWg.Add(1)
	p.tasks <- task
}

// Выключение воркера
func (p *Pool) StopWorker() {
	p.workersWg.Add(1)
	p.quit <- struct{}{}
}

func (p *Pool) startTask() {
	atomic.AddInt32(&p.freeWorkerCount, -1)
	atomic.AddInt32(&p.busyWorkerCount, 1)
	atomic.AddInt32(&p.waitingTaskCount, 1)
}

func (p *Pool) endTask() {
	atomic.AddInt32(&p.busyWorkerCount, -1)
	atomic.AddInt32(&p.freeWorkerCount, 1)
	atomic.AddInt32(&p.succesfullTaskCount, 1)
	atomic.AddInt32(&p.waitingTaskCount, -1)
}
