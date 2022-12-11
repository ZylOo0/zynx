package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(IRequest)
}
