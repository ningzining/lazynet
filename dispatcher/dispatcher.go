package dispatcher

import (
	"math/rand"

	log "github.com/ningzining/L-log"
	"github.com/ningzining/lazynet/iface"
)

type Dispatcher struct {
	workerPoolSize int                  // 工作池大小
	taskQueueSize  int                  // 任务队列长度
	taskQueue      []chan iface.Request // 任务列表
}

func NewDispatcher(workerPoolSize int, taskQueueSize int) iface.Dispatcher {
	return &Dispatcher{
		workerPoolSize: workerPoolSize,
		taskQueueSize:  taskQueueSize,
		taskQueue:      make([]chan iface.Request, workerPoolSize),
	}
}

func (d *Dispatcher) Dispatch(req iface.Request) {
	// 随机分配一个worker
	workerId := rand.Intn(d.workerPoolSize)

	// 写入工作队列
	d.taskQueue[workerId] <- req
}

func (d *Dispatcher) StartWorkerPool() {
	for i := 0; i < d.workerPoolSize; i++ {
		d.taskQueue[i] = make(chan iface.Request, d.taskQueueSize)

		go d.startWorker(i, d.taskQueue[i])
	}
}

func (d *Dispatcher) startWorker(workerId int, taskQueue chan iface.Request) {
	for {
		select {
		case msg := <-taskQueue:
			d.doHandler(workerId, msg)
		}
	}
}

func (d *Dispatcher) doHandler(workerId int, req iface.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("workerId: %d handle msg panic: %v", workerId, err)
		}
	}()

	// 获取请求的连接管道，并处理请求消息
	req.GetConn().GetPipeline().Handle(req.GetMsg())
}
