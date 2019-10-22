package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

//connection 属性
type Connection struct {
	//当前连接的socket tcp套接字
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接的关闭状态
	isclosed bool
	//该链接的处理方法API
	handleAPI ziface.HandFunc

	//告知该连接已经退出/停止的channel

	EXitBuffchan chan bool

	//该连接的处理方法router
	Router ziface.IRouter
}

//处理从conn 读数据

func (c *Connection) StartReader() {

	for {
		//切片
		buf := make([]byte, 512)

		_, err := c.Conn.Read(buf)

		fmt.Println("bufbufbufbufbufbufbufbufbufbuf")
		if err != nil {
			fmt.Println("recv buf err", err)
			c.EXitBuffchan <- true
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	c.EXitBuffchan <- true
		// 	return
		// }
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}

}

//启动链接
func (c *Connection) Start() {

	go c.StartReader()
	for {

		select {

		case <-c.EXitBuffchan:
			return
		}

	}

}

//停止链接
func (c *Connection) Stop() {

	if c.isclosed == true {
		return
	}
	c.isclosed = true

	c.Conn.Close()

	c.EXitBuffchan <- true
	close(c.EXitBuffchan)

}

//从当前连接获取原始的socket tcpconn
func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn

}

//获取当前连接id
func (c *Connection) GetConnID() uint32 {

	return c.ConnID

}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()

}

//创建连接的方法

func NewConntion(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {

	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isclosed:     false,
		Router:       router,
		EXitBuffchan: make(chan bool, 1),
	}
	return c
}
