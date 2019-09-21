package beanstalkg

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const minBuffer  = 1500

type Connection struct {
	Host string
	Port int
	connection *net.TCPConn
	bufRead *bufio.Reader
	bufWrite *bufio.Writer
}

func (c *Connection) getAddr() string {
	return fmt.Sprintf(`%s:%d`, c.Host, c.Port)
}

func dialConnection(c *Connection) (*net.TCPConn, error) {
	conn, err := net.Dial("tcp", c.getAddr())
	if err != nil {
		return nil, err
	}
	c.bufRead = bufio.NewReader(conn)
	c.bufWrite = bufio.NewWriter(conn)
	tcpConn, _ := conn.(*net.TCPConn)
	return tcpConn, nil
}

func NewConnection(host string, port int) (*Connection, error) {
	conn := &Connection{Host: host, Port: port}
	connect, err := dialConnection(conn)
	if err != nil {
		return nil, err
	}
	conn.connection = connect
	fmt.Println(connect)
	return conn, nil
}

func (c *Connection) SendAll(msg []byte) (int, error) {
	written := 0
	forWrite := msg
	var n int
	var err error
	for written < len(msg) {
		forBuff := len(forWrite) >= minBuffer
		if forBuff {
			n, err = sendAllBuffer(c, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		} else {
			n, err = sendAllNoBuffer(c, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		}
		written += n
		forWrite = forWrite[n:]
	}
	return written, nil
}

func sendAllNoBuffer(c *Connection, msg []byte) (int, error) {
	n, err := c.connection.Write(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func sendAllBuffer(c *Connection, msg []byte) (int, error) {
	n, err := c.bufWrite.Write(msg)
	if err != nil {
		return n, err
	}
	err = c.bufWrite.Flush()
	if err != nil {
		return n, err
	}
	return n, nil
}


func (c *Connection) GetResp(cmd string) (string, error) {
	_, err := c.SendAll([]byte(cmd))
	if err != nil {
		return "", err
	}
	resp, err := c.bufRead.ReadString('\n')
	if err != nil {
		return "", err
	}
	return resp, nil
}


func (c *Connection) readBody(msgLen int) ([]byte, error) {
	msgLen += 2
	body := make([]byte, msgLen)
	n, err := io.ReadFull(c.bufRead, body)
	if err != nil {
		return nil, err
	}
	return body[:n-2], nil
}

func isNetTempErr(err error) bool {
	if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
		return true
	}
	return false
}