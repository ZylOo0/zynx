package znet

type Handler interface {
	PreHandle(request *Request)
	InHandle(request *Request)
	PostHandle(request *Request)
}

//type HandlerBase struct{}
//
//func (h *HandlerBase) PreHandle(request *Request) {}
//
//func (h *HandlerBase) InHandle(request *Request) {}
//
//func (h *HandlerBase) PostHandle(request *Request) {}
