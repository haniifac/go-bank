package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	userAgentHeader        = "user-agent"
	ipAddressHeader        = "x-forwarded-for"
	gatewayUserAgentHeader = "grpcgateway-user-agent"
)

type Metadata struct {
	UserAgent string
	IpAddress string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}

	if m, ok := metadata.FromIncomingContext(ctx); ok {
		// log.Printf("Metadata context: %+v\n", m)

		if userAgent := m.Get(userAgentHeader); len(userAgent) > 0 {
			md.UserAgent = userAgent[0]
		}

		if gatewayUserAgent := m.Get(gatewayUserAgentHeader); len(gatewayUserAgent) > 0 {
			md.UserAgent = gatewayUserAgent[0]
		}

		if ipAddress := m.Get(ipAddressHeader); len(ipAddress) > 0 {
			md.IpAddress = ipAddress[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		md.IpAddress = p.Addr.String()
	}

	return md
}
