package znet

import "zyloo.com/zinx/ziface"

type Request struct {
	conn ziface.IConnection // 已经与客户端建立好的连接
	msg  ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
