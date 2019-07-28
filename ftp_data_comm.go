package main

import (
	"net"
	"sync"
)

type FTPDataComm struct {
	LocalConnection  net.Conn
	RemoteConnection net.Conn
	Napter           *PortMapping
	Soracom          *SoracomAdapter
	Closed           bool
	Mutex            sync.Mutex
}

func (comm *FTPDataComm) Start() {
	comm.Closed = false
	go comm.ProcessLocalReceiveLoop()
	go comm.ProcessRemoteReceiveLoop()
}

func (comm *FTPDataComm) ProcessLocalReceiveLoop() {
	buf := make([]byte, 65536)
	for {
		readLen, err := comm.LocalConnection.Read(buf)
		if err != nil && err.Error() != "EOF" {
			break
		}
		comm.RemoteConnection.Write(buf[:readLen])
		if err != nil && err.Error() == "EOF" {
			break
		}
	}
	comm.Close()
}

func (comm *FTPDataComm) ProcessRemoteReceiveLoop() {
	buf := make([]byte, 65536)
	for {
		readLen, err := comm.RemoteConnection.Read(buf)
		if err != nil && err.Error() != "EOF" {
			break
		}
		comm.LocalConnection.Write(buf[:readLen])
		if err != nil && err.Error() == "EOF" {
			break
		}
	}
	comm.Close()
}

func (comm *FTPDataComm) Close() {
	comm.Mutex.Lock()
	defer comm.Mutex.Unlock()
	if !comm.Closed {
		comm.RemoteConnection.Close()
		comm.LocalConnection.Close()
		comm.Soracom.StopNapter(comm.Napter)
	}
	comm.Closed = true
}
