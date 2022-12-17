# zynx

原项目名 zinx，是刘丹冰大神写的 TCP 长连接的服务端框架，zynx 在此基础上进行修改。

## TODO

- [x] 调整 router 和 handler 的名称
- [x] 去除不必要的接口
- [x] 将 GlobalObject 改为 config
- [x] 删除 datapack 类，改为函数
- [x] 调整 msg 和 data 的名称
- [x] 封装 Connection 的读操作 ReadMsg
- [x] 优化工作池负载均衡，改用 Request ID
- [ ] 读写分离有什么好处？
- [ ] 将不必要的大写改为小写