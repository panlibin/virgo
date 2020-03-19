package virgo

import (
	"syscall"
	"os/signal"
	"runtime"
	"github.com/panlibin/vglog"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	taskTypeNew int32 = iota
	taskTypeFinish
	taskTypeQuit
)

const _DefaultMainChannelSize = 1024

type task struct {
	taskType int32
	f func([]interface{})
	args []interface{}
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
	wg sync.WaitGroup
	mainChan chan *task
	pending int32
	quitFlag int32
	service IService
}

// NewProcedure 创建
func NewProcedure(s IService) *Procedure {
	p := &Procedure{
		mainChan: make(chan *task, _DefaultMainChannelSize),
		service: s,
	}
	return p
}

// SyncTask 主线程内执行函数
func (p *Procedure) SyncTask(f func([]interface{}), args ...interface{}) {
	atomic.AddInt32(&p.pending, 1)
	p.mainChan <- &task{
		f: f,
		args: args,
		taskType: taskTypeNew,
	}
}

// AsyncTask 主线程外执行函数,不直接go,为了统计运行中的任务
func (p *Procedure) AsyncTask(f func([]interface{}), args ...interface{}) {
	atomic.AddInt32(&p.pending, 1)
	go func ()  {
		protectedExecute(f, args)
		p.mainChan <- &task{taskType: taskTypeFinish}
	}()
}

// AfterFunc 主线程内定时回调
func (p *Procedure) AfterFunc(d time.Duration, f func([]interface{}), args ...interface{}) *time.Timer {
	return time.AfterFunc(d, func ()  {
		p.SyncTask(f, args...)
	})
}

// Start 启动
func (p *Procedure) Start() {
	p.run()
	p.SyncTask(func([]interface{}){
		p.service.OnInit(p)
	})
}

// Stop 停止
func (p *Procedure) Stop() {
	atomic.AddInt32(&p.pending, 1)
	go func ()  {
		p.mainChan <- &task{
			f: func([]interface{}){
				p.service.OnRelease()
			},
			taskType: taskTypeQuit,
		}
	}()
}

func (p *Procedure) run() {
	p.wg.Add(1)
	go func() {
		for pTask := range p.mainChan {
			switch pTask.taskType {
			case taskTypeNew:
				protectedExecute(pTask.f, pTask.args)
			case taskTypeQuit:
				protectedExecute(pTask.f, pTask.args)
				p.quitFlag = 1
			}

			tmpPending := atomic.AddInt32(&p.pending, -1)
			if p.quitFlag == 1 && tmpPending <= 0 {
				break
			}
		}

		close(p.mainChan)
		p.mainChan = nil
		p.wg.Done()
	}()
}

func (p *Procedure) waitQuit() {
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		bQuit := false
		for !bQuit {
			sig := <- sigChan
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
