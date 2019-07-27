package main

type PortMapping struct {
	IpAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
	Hostname  string `json:"hostname"`
}
