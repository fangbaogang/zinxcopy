package ziface

type IServer interface {

	//启动服务器的方法

	Start()

	//停止服务器方法
	Stop()

	//开启业务服务方法
	Server()

	//路由功能
	AddRouter(router IRouter)
}
