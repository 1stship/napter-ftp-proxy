package main

import (
	"crypto/tls"
	"net"
	"strconv"
)

type FTPDataListener struct {
	LocalAddress string
	LocalPort    int
	Insecure     bool
	Napter       *PortMapping
	Soracom      *SoracomAdapter
	Listener     net.Listener
}

func (listener *FTPDataListener) Start() error {
	listen, err := net.Listen("tcp", listener.LocalAddress+":"+strconv.Itoa(listener.LocalPort))
	if err != nil {
		return err
	}
	listener.Listener = listen
	go func() {
		for {
			localConnection, err := listen.Accept()
			if err != nil {
				break
			}
			var remoteConnection net.Conn
			if listener.Insecure {
				remoteConnection, err = net.Dial("tcp", listener.Napter.Hostname+":"+strconv.Itoa(listener.Napter.Port))
			} else {
				tlsConf := &tls.Config{MinVersion: tls.VersionTLS12}
				remoteConnection, err = tls.Dial("tcp", listener.Napter.Hostname+":"+strconv.Itoa(listener.Napter.Port), tlsConf)
			}

			if err != nil {
				localConnection.Close()
				break
			}
			ftpDataComm := &FTPDataComm{
				LocalConnection:  localConnection,
				RemoteConnection: remoteConnection,
				Napter:           listener.Napter,
				Soracom:          listener.Soracom}
			ftpDataComm.Start()
		}
	}()
	return nil
}

func (listener *FTPDataListener) Stop() {
	listener.Listener.Close()
}
