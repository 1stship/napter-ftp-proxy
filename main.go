package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	const version = "0.1.0"
	dispVersion := false
	var target string
	var listenAddress string
	var localPort int
	var remotePort int
	var insecure bool
	flag.BoolVar(&dispVersion, "v", false, "バージョン表示")
	flag.BoolVar(&dispVersion, "version", false, "バージョン表示")
	flag.StringVar(&listenAddress, "listen", "127.0.0.1", "待ち受けIPアドレス")
	flag.StringVar(&target, "target", "", "接続先のIMSI")
	flag.BoolVar(&insecure, "insecure", false, "暗号化しない")
	flag.IntVar(&localPort, "local", 21, "ローカル待ち受けポート")
	flag.IntVar(&remotePort, "remote", 21, "接続先ポート")
	flag.Parse()

	if dispVersion {
		fmt.Printf("napter-ftp-proxy: Ver %s\n", version)
		os.Exit(0)
	}

	if target == "" {
		fmt.Println("--target で接続先のIMSIを指定してください")
		os.Exit(1)
	}

	email := os.Getenv("SORACOM_EMAIL")
	password := os.Getenv("SORACOM_PASSWORD")

	if email == "" {
		email = getInput("Input Soracom account email: ")
	}

	if password == "" {
		password = getPasswordInput("Input Soracom account password: ")
	}

	credential := &SoracomCredential{Email: email, Password: password}
	soracom := &SoracomAdapter{credential: credential, target: target}
	soracom.GetSoracomToken()
	destination := &PortMappingDestination{Imsi: target, Port: remotePort}
	createRequest := &CreatePortMappingRequest{Destination: destination, TlsRequired: !insecure}
	portMapping, err := soracom.StartNapter(createRequest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fail to start napter")
		os.Exit(1)
	}
	defer func() {
		err = soracom.StopNapter(portMapping)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fail to stop napter")
			os.Exit(1)
		}
	}()

	commandListener := &FTPCommandListener{
		LocalAddress: listenAddress,
		LocalPort:    localPort,
		Napter:       portMapping,
		Soracom:      soracom,
		Insecure:     insecure}

	err = commandListener.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	trapSignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, trapSignals...)

	<-sigCh
}

func getInput(inst string) string {
	for {
		fmt.Print(inst)
		scanner := bufio.NewScanner(os.Stdin)
		done := scanner.Scan()
		if done {
			input := scanner.Text()
			if input != "" {
				return input
			}
		}
	}
}

func getPasswordInput(inst string) string {
	for {
		fmt.Print(inst)
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			continue
		} else {
			if string(password) != "" {
				fmt.Println("")
				return string(password)
			}
		}
	}
}
