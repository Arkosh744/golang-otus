package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type taskManagement struct {
	m               int
	errCount        int32
	breakOnErrLimit bool
	wg              sync.WaitGroup
}

func (te *taskManagement) incErrCount() {
	atomic.AddInt32(&te.errCount, 1)
}

func (te *taskManagement) getErrCount() int32 {
	return atomic.LoadInt32(&te.errCount)
}

func (te *taskManagement) checkErrorsLimitEquals() bool {
	return te.getErrCount() == int32(te.m) && te.breakOnErrLimit
}

func (te *taskManagement) checkErrorsLimitExceeded() bool {
	return te.getErrCount() >= int32(te.m) && te.breakOnErrLimit
}

func newTaskManagement(m int) *taskManagement {
	te := taskManagement{m: m}
	if m > 0 {
		te.breakOnErrLimit = true
	}
	return &te
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksMan := newTaskManagement(m)

	tasksMan.wg.Add(n)

	cancel := make(chan struct{})
	defer close(cancel)

	taskChan := taskWriter(cancel, tasks)

	for i := 0; i < n; i++ {
		go taskWorker(cancel, taskChan, tasksMan)
	}

	tasksMan.wg.Wait()

	if tasksMan.errCount >= int32(tasksMan.m) && tasksMan.breakOnErrLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func taskWorker(cancel chan struct{}, taskChan <-chan Task, tasksMan *taskManagement) {
	defer tasksMan.wg.Done()

	for {
		if tasksMan.checkErrorsLimitEquals() {
			cancel <- struct{}{}
			return
		}

		select {
		case task, ok := <-taskChan:
			if !ok || tasksMan.checkErrorsLimitExceeded() {
				return
			}

			if err := task(); err != nil {
				tasksMan.incErrCount()
			}

		case <-cancel:
			return
		}
	}
}

func taskWriter(cancel <-chan struct{}, tasks []Task) chan Task {
	taskChan := make(chan Task)
	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-cancel:
				return
			}
		}
	}()
	return taskChan
}
