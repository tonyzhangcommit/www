package bussinesscode

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

/*
	实现协程池
		why: 原生协程goroutine启动方便，可以方便支持高效并发，但是goroutine太多会导致调度性能下降，GC频繁，内存暴涨等问题；因此需要这样一个场景：
		限制协程数量。
		Q&A:
			1. 怎么限制goroutine的数量？
			2. goroutine怎么重用？
			3. 任务如何执行？
		模型：
			生产者消费者模型
			这里golang 协程是使用关键字调用，所以不能用连接池的思路来解决
*/

/*
	定义任务：
		任务中需要有执行的函数和对应的参数，但是参数类型个数都不确定，需要解决下
*/

type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

/*
定义任务池，要考虑下面几个参数：容量，当前worker的数量，任务队列，任务池状态
*/
type Pool struct {
	capacity       uint64
	runningworkers uint64
	status         int64
	chTask         chan *Task
	sync.Mutex                       // 互斥锁
	PanicHandler   func(interface{}) // 异常处理函数
}

// 任务池初始函数
var ErrorInvalidPoolCap = errors.New("非法任务池容量")

const (
	RUNNING = 1 // 运行中
	STOPED  = 0 // 停止
)

func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrorInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		chTask:   make(chan *Task, capacity),
	}, nil
}

// 启动worker
// 由于这里涉及到并发操作，++ 操作会有数据竞争， 这里需要更改为原子操作
func (p *Pool) Run() {
	p.inCreRunning()
	go func() {
		defer func() {
			// worker 结束， 释放
			p.deCreRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("worker panic:%s\n", r)
				}
			}
			p.checkworker()
		}()
		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					// 队列中已经没有待处理的任务了
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

// 原子操作runnings
func (p *Pool) inCreRunning() {
	atomic.AddUint64(&p.runningworkers, 1)
}

func (p *Pool) deCreRunning() {
	atomic.AddUint64(&p.runningworkers, ^uint64(0))
}

func (p *Pool) getRunning() uint64 {
	return atomic.LoadUint64(&p.runningworkers)
}

// 获取任务池容量，因为这个在初始化时固定容量，不需要考虑并发
func (p *Pool) getCap() uint64 {
	return p.capacity
}

// 对status 状态进行加锁操作
func (p *Pool) setStatus(status int64) bool {
	p.Lock()
	defer p.Unlock()

	if p.status == status {
		return false
	}
	p.status = status
	return true
}

// 任务入池
func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()

	if p.status == STOPED {
		return ErrorPoolClosed
	}

	if p.getRunning() < p.getCap() {
		p.Run()
	}
	if p.status == RUNNING {
		p.chTask <- task
	}
	return nil
}

// 关闭任务池
var ErrorPoolClosed = errors.New("协程池已关闭")

func (p *Pool) Close() {
	p.setStatus(STOPED)
	// 消费完当前chan 中所有的task
	for len(p.chTask) > 0 {
		time.Sleep(time.Microsecond * 100)
	}
	close(p.chTask)
}

func (p *Pool) checkworker() {
	p.Lock()
	defer p.Unlock()
	if p.runningworkers == 0 && len(p.chTask) > 0 {
		p.Run()
	}
}
