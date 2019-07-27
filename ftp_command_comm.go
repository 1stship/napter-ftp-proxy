package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type FTPCommandComm struct {
	LocalAddress     string
	Insecure         bool
	LocalConnection  net.Conn
	RemoteConnection net.Conn
	Soracom          *SoracomAdapter
	Listeners        []*FTPDataListener
	Closed           bool
	Mutex            sync.Mutex
}

func (comm *FTPCommandComm) Start() {
	fmt.Println("start ftp connection")
	comm.Closed = false
	go comm.ProcessLocalReceiveLoop()
	go comm.ProcessRemoteReceiveLoop()
}

func (comm *FTPCommandComm) ProcessLocalReceiveLoop() {
	scanner := bufio.NewScanner(comm.LocalConnection)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		fmt.Fprintln(comm.RemoteConnection, scanner.Text())
	}
	comm.Close()
}

func (comm *FTPCommandComm) ProcessRemoteReceiveLoop() {
	scanner := bufio.NewScanner(comm.RemoteConnection)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		processedText := comm.ProcessCommand(scanner.Text())
		fmt.Fprintln(comm.LocalConnection, processedText)
	}
	comm.Close()
}

func (comm *FTPCommandComm) ProcessCommand(command string) string {
	ret := command
	if strings.HasPrefix(command, "227") {
		leftParenthesisIndex := strings.Index(command, "(")
		rightParenthesisIndex := strings.Index(command, ")")
		passiveParamsString := command[leftParenthesisIndex+1 : rightParenthesisIndex]
		passiveParams := strings.Split(passiveParamsString, ",")
		passivePortUpper, _ := strconv.Atoi(passiveParams[4])
		passivePortLower, _ := strconv.Atoi(passiveParams[5])
		passivePort := passivePortUpper*256 + passivePortLower
		destination := &PortMappingDestination{Imsi: comm.Soracom.target, Port: passivePort}
		createRequest := &CreatePortMappingRequest{Destination: destination, TlsRequired: !comm.Insecure}
		portMapping, _ := comm.Soracom.StartNapter(createRequest)
		dataListner := &FTPDataListener{
			LocalAddress: comm.LocalAddress,
			LocalPort:    passivePort,
			Insecure:     comm.Insecure,
			Napter:       portMapping,
			Soracom:      comm.Soracom}
		dataListner.Start()
		comm.Listeners = append(comm.Listeners, dataListner)
		ret = "227 Entering Passive Mode (" + strings.Replace(comm.LocalAddress, ".", ",", -1) + "," + passiveParams[4] + "," + passiveParams[5] + ")."
	} else if strings.HasPrefix(command, "226") {
		for _, dataListener := range comm.Listeners {
			dataListener.Stop()
		}
		comm.Listeners = make([]*FTPDataListener, 0)
	}
	return ret
}

func (comm *FTPCommandComm) Close() {
	comm.Mutex.Lock()
	defer comm.Mutex.Unlock()
	if !comm.Closed {
		fmt.Println("finish ftp connection")
		comm.LocalConnection.Close()
		comm.RemoteConnection.Close()
	}
	comm.Closed = true
}
