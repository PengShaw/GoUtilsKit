package socket

import (
	"io"
	"net"

	"github.com/PengShaw/GoUtilsKit/logger"
)

// RunSocketClient builds a socket connection, and send data to server
func RunSocketClient(network, address string, ch <-chan []byte) {
	conn, err := net.Dial(network, address)
	if err != nil {
		logger.Errorf("connect to %s:%s failed: %s", network, address, err)
		return
	}
	defer conn.Close()
	logger.Infof("dial: <%s>", conn.RemoteAddr().String())

	for {
		data := <-ch
		_, err := conn.Write(data)
		if err != nil {
			logger.Errorf("send data to %s:%s failed: %s", network, address, err)
			logger.Debugf("send data to %s:%s failed: %s: %s", network, address, err, data)
		}
	}
}

// RunUDPServer listens an udp socket, and send received data to channel
func RunUDPServer(address string, ch chan<- []byte) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		logger.Errorf("listen udp:%s failed: %s", address, err)
		return
	}
	defer conn.Close()
	logger.Infof("listen: <%s>", conn.LocalAddr().String())

	buf := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			logger.Errorf("listen udp:%s data failed: %s", address, err)
			continue
		}
		logger.Infof("received data from %s", addr.String())
		logger.Debugf("received data from %s: %s", addr.String(), buf)
		ch <- buf
	}
}

// RunTCPServer listens an tcp socket, and send received data to channel
func RunTCPServer(address string, ch chan<- []byte) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Errorf("listen tcp:%s failed: %s", address, err)
		return
	}
	defer l.Close()
	logger.Debug("tcp listen build")

	for {
		conn, err := l.Accept()
		logger.Infof("listen: <%s>", l.Addr().String())
		if err != nil {
			logger.Errorf("connect tcp:%s failed: %s", address, err)
			continue
		}

		go func(c net.Conn, ch chan<- []byte) {
			defer c.Close()
			for {
				buf := make([]byte, 1024)
				_, err := c.Read(buf)
				if err != nil && err != io.EOF {
					logger.Errorf("listen tcp:%s data failed: %s", address, err)
					break
				}
				if err == io.EOF {
					break
				}
				logger.Infof("received data from %s", c.RemoteAddr().String())
				logger.Debugf("received data from %s: %s", c.RemoteAddr().String(), buf)
				ch <- buf
			}
		}(conn, ch)
	}
}

// RunUnixServer listens an unix domain socket, and send received data to channel
func RunUnixServer(address string, ch chan<- []byte) {
	l, err := net.Listen("unix", address)
	if err != nil {
		logger.Errorf("listen unix:%s failed: %s", address, err)
		return
	}
	defer l.Close()
	logger.Debug("unix listen build")

	for {
		conn, err := l.Accept()
		logger.Infof("listen: <%s>", l.Addr().String())
		if err != nil {
			logger.Errorf("connect unix:%s failed: %s", address, err)
			continue
		}

		go func(c net.Conn, ch chan<- []byte) {
			defer c.Close()
			for {
				buf := make([]byte, 1024)
				_, err := c.Read(buf)
				if err != nil && err != io.EOF {
					logger.Errorf("listen unix:%s data failed: %s", address, err)
					break
				}
				if err == io.EOF {
					break
				}
				logger.Infof("received data from %s", c.LocalAddr().String())
				logger.Debugf("received data from %s: %s", c.LocalAddr().String(), buf)
				ch <- buf
			}
		}(conn, ch)
	}
}
