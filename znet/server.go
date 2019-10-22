package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/ziface"
)

//Server is a struct type
type Server struct {

	//服务器名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务器绑定的IP地址
	IP string
	//服务器端口
	Port uint32
	//当前server由用户绑定的回调router,也就是server注册的连接对应的处理业务
	Router ziface.IRouter
}

//--------------------实现ziface.iserver接口里面的所有方法-------------------------

//  服务器开始
func (s *Server) Start() {

	fmt.Printf("[START] ZINX server listen at IP :%s,port %d,is starting \n", s.IP, s.Port)

	//开启listen 业务
	go func() {
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {

		}

		//监听服务器
		tcplistener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {

		}
		var cid uint32
		cid = 0
		//3	启动server网络连接业务
		for {
			//3.1	阻塞等待客户端建立连接请求
			conn, err := tcplistener.AcceptTCP()
			if err != nil {
				continue
			}

			dealconn := NewConntion(conn, cid, s.Router)
			cid++

			go dealconn.Start()
		}

	}()

}

//Stop is a function.
func (s *Server) Stop() {

	fmt.Println("[STOP] zinx server ,name \n", s.Name)
}

//服务器
func (s *Server) Server() {

	s.Start()

	//阻塞？到底是什么意思？
	for {
		time.Sleep(time.Second)
	}
}

//为什么返回的是一个接口类型
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp",
		IP:        "127.0.0.1",
		Port:      7777,
		Router:    nil,
	}
	return s
}

//定义当前客户端连接的handle api

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	//回显业务
	fmt.Println("回显")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("回显失败")
		return errors.New("cccc")
	}
	return nil
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router

	fmt.Println("add router succ!")

}
