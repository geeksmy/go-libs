package grpcclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var ClientOptions = []grpc.DialOption{
	grpc.WithInsecure(),
	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(30) * time.Second,
		Timeout:             time.Duration(10) * time.Second,
		PermitWithoutStream: true,
	}),
}

var clientCache sync.Map

func NewClient(host string, port int) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)
	if conn, ok := clientCache.Load(endpoint); ok {
		return conn.(*grpc.ClientConn), nil
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, endpoint, ClientOptions...)

	if err != nil {
		return nil, err
	}

	clientCache.Store(endpoint, conn)

	return conn, nil
}

func NewTLSClient(host string, port int, tlsConfig *tls.Config) (*grpc.ClientConn, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)
	if conn, ok := clientCache.Load(endpoint); ok {
		return conn.(*grpc.ClientConn), nil
	}
	creds := credentials.NewTLS(tlsConfig)
	tlsClientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(30) * time.Second,
			Timeout:             time.Duration(10) * time.Second,
			PermitWithoutStream: true,
		}),
	}
	conn, err := grpc.Dial(endpoint, tlsClientOptions...)
	if err != nil {
		return nil, err
	}
	clientCache.Store(endpoint, conn)
	return conn, nil
}

func NewClientWithEndpoint(endpoint string) (*grpc.ClientConn, error) {
	host, port, err := ParseToHostPort(endpoint)
	if err != nil {
		return nil, err
	}

	return NewClient(host, port)
}

func ParseToHostPort(endpoint string) (host string, port int, err error) {
	if strings.HasPrefix(endpoint, "grpc://") || strings.HasPrefix(endpoint, "grpcs://") {
		URL, err1 := url.Parse(endpoint)
		if err1 != nil {
			err = err1
			return
		}
		endpoint = URL.Host
	}

	args := strings.Split(endpoint, ":")
	if len(args) != 2 {
		err = fmt.Errorf("invalid server host: %s", endpoint)
		return
	}
	host = args[0]
	port, err = strconv.Atoi(args[1])
	return
}
