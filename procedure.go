package virgo

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	logger "github.com/panlibin/vglog"
)

const (
	taskTypeNew int32 = iota
	taskTypeFinish
	taskTypeQuit
)

const _DefaultMainChannelSize = 1024

type task struct {
	taskType int32
	f        func([]interface{})
	args     []interface{}
}

// IProcedure 主线程接口
type IProcedure interface {
	SyncTask(f func([]interface{}), args ...interface{})
	AsyncTask(f func([]interface{}), args ...interface{})
	AfterFunc(d time.Duration, f func([]interface{}), args ...interface{}) *time.Timer
	Start()
	Stop()
}

// Procedure 主线程
type Procedure struct {
	wg         sync.WaitGroup
	mainChan   chan *task
	pending    int32
	quitFlag   int32
	service    IService
	runningQue *taskQueue
	pendingQue *taskQueue
	cond       *sync.Cond
}

// NewProcedure 创建
func NewProcedure(s IService) *Procedure {
	p := &Procedure{
		runningQue: newTaskQueue(_DefaultMainChannelSize),
		pendingQue: newTaskQueue(_DefaultMainChannelSize),
		service:    s,
		cond:       sync.NewCond(&sync.Mutex{}),
	}
	return p
}

// SyncTask 主线程内执行函数
func (p *Procedure) SyncTask(f func([]interface{}), args ...interface{}) {
	atomic.AddInt32(&p.pending, 1)

	p.pushTask(&task{
		f:        f,
		args:     args,
		taskType: taskTypeNew,
	})
}

// AsyncTask 主线程外执行函数,不直接go,为了统计运行中的任务
func (p *Procedure) AsyncTask(f func([]interface{}), args ...interface{}) {
	atomic.AddInt32(&p.pending, 1)
	go func() {
		protectedExecute(f, args)

		p.pushTask(&task{
			taskType: taskTypeFinish,
		})
	}()
}

// AfterFunc 主线程内定时回调
func (p *Procedure) AfterFunc(d time.Duration, f func([]interface{}), args ...interface{}) *time.Timer {
	return time.AfterFunc(d, func() {
		p.SyncTask(f, args...)
	})
}

// Start 启动
func (p *Procedure) Start() {
	p.run()
	p.SyncTask(func([]interface{}) {
		p.service.OnInit(p)
	})
}

// Stop 停止
func (p *Procedure) Stop() {
	atomic.AddInt32(&p.pending, 1)

	p.pushTask(&task{
		f: func([]interface{}) {
			p.service.OnRelease()
		},
		taskType: taskTypeQuit,
	})
}

func (p *Procedure) run() {
	p.wg.Add(1)
	go func() {
		var tmpPending int32
		var pTask *task
		for {
			p.cond.L.Lock()
			if p.pendingQue.empty {
				p.cond.Wait()
				p.cond.L.Unlock()
			} else {
				p.pendingQue, p.runningQue = p.runningQue, p.pendingQue
				p.cond.L.Unlock()

				for pTask = p.runningQue.pop(); pTask != nil; pTask = p.runningQue.pop() {
					switch pTask.taskType {
					case taskTypeNew:
						protectedExecute(pTask.f, pTask.args)
					case taskTypeQuit:
						protectedExecute(pTask.f, pTask.args)
						p.quitFlag = 1
					}
					tmpPending = atomic.AddInt32(&p.pending, -1)
				}

				if p.quitFlag == 1 && tmpPending <= 0 {
					break
				}
			}
		}

		p.wg.Done()
	}()
}

func (p *Procedure) pushTask(pTask *task) {
	p.cond.L.Lock()
	p.pendingQue.push(pTask)
	p.cond.L.Unlock()
	p.cond.Signal()
}

func (p *Procedure) waitQuit() {
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		bQuit := false
		for !bQuit {
			sig := <-sigChan
			if sig == os.Signal(syscall.SIGINT) || sig == os.Signal(syscall.SIGTERM) {
				bQuit = true
			}
		}
		p.Stop()
	}()

	p.wg.Wait()
}

func protectedExecute(f func([]interface{}), args []interface{}) (err interface{}) {
	defer func() {
		if err = recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.Errorf("%v\n%s", err, buf[:n])
		}
	}()

	f(args)

	return
}
