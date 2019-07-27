package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
)

type FTPCommandListener struct {
	LocalAddress string
	LocalPort    int
	Insecure     bool
	Napter       *PortMapping
	Soracom      *SoracomAdapter
}

func (listener *FTPCommandListener) Start() error {
	listen, err := net.Listen("tcp", listener.LocalAddress+":"+strconv.Itoa(listener.LocalPort))
	if err != nil {
		return err
	}
	fmt.Println("start ftp listen")
	go func() {
		for {
			localConnection, err := listen.Accept()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			var remoteConnection net.Conn
			if listener.Insecure {
				remoteConnection, err = net.Dial("tcp", listener.Napter.Hostname+":"+strconv.Itoa(listener.Napter.Port))
			} else {
				tlsConf := &tls.Config{MinVersion: tls.VersionTLS12}
				remoteConnection, err = tls.Dial("tcp", listener.Napter.Hostname+":"+strconv.Itoa(listener.Napter.Port), tlsConf)
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				localConnection.Close()
				continue
			}
			ftpCommandComm := &FTPCommandComm{
				LocalAddress:     listener.LocalAddress,
				Insecure:         listener.Insecure,
				LocalConnection:  localConnection,
				RemoteConnection: remoteConnection,
				Soracom:          listener.Soracom,
				Listeners:        make([]*FTPDataListener, 0)}
			ftpCommandComm.Start()
		}
	}()
	return nil
}
