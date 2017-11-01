package main

import (
    "math/rand"
    "time"
    "errors"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

type workerFunc func()

const (
    DEFAULT_WORKERS = 5
)

type WorkerPool struct {
    workers   []*worker
    closeChan chan struct{}
}

func NewWorkerPool(vol int) *WorkerPool {
    if vol <= 0 {
        vol = DEFAULT_WORKERS
    }

    pool := &WorkerPool{
        workers:   make([]*worker, vol),
        closeChan: make(chan struct{}),
    }

    for i, _ := range pool.workers {
        pool.workers[i] = newWorker(i, 1024, pool.closeChan)
        if pool.workers[i] == nil {
            panic("failed to create worker")
        }
    }

    return pool
}

func (wp *WorkerPool) Put(k interface{}, cb func()) error {
    return wp.workers[rand.Intn(10000) % len(wp.workers)].put(workerFunc(cb))
}

func (wp *WorkerPool) Close() {
    close(wp.closeChan)
}

type worker struct {
    index        int
    callbackChan chan workerFunc
    closeChan    chan struct{}
}

func newWorker(i int, c int, closeChan chan struct{}) *worker {
    w := &worker{
        index:        i,
        callbackChan: make(chan workerFunc, c),
        closeChan:    closeChan,
    }
    go w.start()
    return w
}

func (w *worker) start() {
    defer close(w.callbackChan)

    for {
        select {
        case <-w.closeChan:
            break
        case cb := <-w.callbackChan:
            cb()
        }
    }
}

func (w *worker) put(cb workerFunc) error {
    select {
    case w.callbackChan <- cb:
        return nil
    default:
        return errors.New("put task would block")
    }
}
