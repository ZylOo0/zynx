package znet

type Request struct {
	conn *Connection // 已经与客户端建立好的连接
	msg  *Message
}

func (r *Request) GetConnection() *Connection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
