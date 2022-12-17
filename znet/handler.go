package znet

type Handler interface {
	PreHandle(request *Request)
	InHandle(request *Request)
	PostHandle(request *Request)
}

//type VHandler struct{}
//
//func (h *VHandler) PreHandle(request *Request) {}
//
//func (h *VHandler) InHandle(request *Request) {}
//
//func (h *VHandler) PostHandle(request *Request) {}
