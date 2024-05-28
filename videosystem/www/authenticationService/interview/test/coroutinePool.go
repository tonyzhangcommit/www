package test

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

/*
Go 使用了基于 M:N 调度模型的调度器，即多个 goroutine（N）可以在多个操作系统线程（M）上调度执行。这种调度模型的设计意图是为了利用多核处理器，同时减少系统调度的开销。但当存在大量 goroutine 时，以下因素可能导致调度性能下降：

- **上下文切换：** 尽管 goroutine 的上下文切换比操作系统线程轻便，但极大数量的 goroutine 仍然可能导致频繁的切换，从而消耗更多的 CPU 时间。
- **调度队列管理：** 调度器需要管理所有可运行的 goroutine 的队列。当这个队列非常长时，维护队列（如添加、删除、查找下一个可运行的 goroutine）的开销将增加。
- **负载均衡：** Go 调度器还需尝试在多个线程（M）间平衡 goroutine 的执行，以充分利用多核优势。在大量 goroutine 的情况下，有效的负载均衡变得更加复杂和耗时。

### 2. GC 频繁的底层原因

Go 的垃圾收集器是并发的、标记-清除（mark-and-sweep）类型。GC 过程主要分为两个阶段：标记阶段和清除阶段。当有大量的 goroutine 运行时，以下因素会影响 GC 频率：

- **堆内存增长：** 每个 goroutine 可能会有自己的局部变量和堆分配，随着 goroutine 数量的增加，整体的内存分配也会增加，这可能导致堆快速增长。
- **GC 触发机制：** Go 的 GC 设计为自适应触发，当堆内存增长超过一定比例后，GC 就会启动。因此，大量的内存分配会导致更频繁的 GC。
- **标记负载：** 在标记阶段，GC 需要遍历所有活动对象并进行标记。大量的 goroutine 可能意味着更多的栈（每个 goroutine 有自己的栈），每个栈中可能引用了许多对象，这增加了标记阶段的工作量。

### 3. 内存暴涨的底层原因

如上所述，每个 goroutine 虽然在初始时只需要很小的栈空间（通常是几 KB），但以下因素可能导致整体内存使用急剧增加：

- **栈大小动态调整：** 如果 goroutine 的执行需要更多的栈空间，Go 运行时会自动增加其栈的大小。在极端情况下，这可能导致单个 goroutine 占用的内存远超初始分配。
- **堆分配：** goroutine 在运行过程中可能会进行大量的堆分配，特别是在使用复杂的数据结构时（如切片、映射和通道等

）。
- **内存泄露：** 如果 goroutine 中的对象因为某些原因（例如闭包引用）未能正确回收，这些对象将继续占用内存，可能导致内存泄露，进一步加剧内存使用的增加。

因此，合理控制 goroutine 的数量和生命周期，以及优化内存使用策略，是提高 Go 程序性能和稳定性的关键。

*/
