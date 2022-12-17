package znet

import (
	"fmt"
	"strconv"
)

// Router 存储各种 Message ID 对应的 Handler
// 并且将 Request 分配给 Worker
type Router struct {
	Handlers       map[uint32]Handler // msgID 对应的处理方法
	TaskQueue      []chan *Request
	WorkerPoolSize uint32
}

func NewRouter() *Router {
	return &Router{
		Handlers:       make(map[uint32]Handler),
		WorkerPoolSize: config.WorkerPoolSize,
		TaskQueue:      make([]chan *Request, config.WorkerPoolSize),
	}
}

func (r *Router) TriggerHandler(request *Request) {
	handler, ok := r.Handlers[request.GetMsgID()]
	if !ok {
		fmt.Println("handler msgID = ", request.GetMsgID(), " is NOT FOUND! Need register!")
		return
	}
	handler.PreHandle(request)
	handler.InHandle(request)
	handler.PostHandle(request)
}

func (r *Router) AddHandler(msgID uint32, handler Handler) {
	if _, ok := r.Handlers[msgID]; ok { // ID 已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	r.Handlers[msgID] = handler
	fmt.Println("Add handler MsgID = ", msgID, "success!")
}

func (r *Router) StartWorkerPool() {
	for i := 0; i < int(r.WorkerPoolSize); i++ {
		r.TaskQueue[i] = make(chan *Request, config.MaxWorkerTaskLen)
		go r.StartOneWorker(i, r.TaskQueue[i])
	}
}

func (r *Router) StartOneWorker(workerID int, taskQueue chan *Request) {
	fmt.Println("Worker ID = ", workerID, "started...")

	for {
		select {
		case request := <-taskQueue:
			r.TriggerHandler(request)
		}
	}
}

func (r *Router) SendRequestToTaskQueue(request *Request) {
	workerID := request.GetReqID() % r.WorkerPoolSize
	fmt.Println("Add ReqID = ", request.GetReqID(),
		"MsgID = ", request.GetMsgID(),
		"to WorkerID = ", workerID)
	r.TaskQueue[workerID] <- request
}
