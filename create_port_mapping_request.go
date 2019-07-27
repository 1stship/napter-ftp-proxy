package main

type CreatePortMappingRequest struct {
	Destination *PortMappingDestination `json:"destination"`
	TlsRequired bool                    `json:"tlsRequired"`
}
