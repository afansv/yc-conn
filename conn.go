package ycconn

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/dial"
	"github.com/yandex-cloud/go-sdk/pkg/grpcclient"
	"github.com/yandex-cloud/go-sdk/pkg/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionConfig struct {
	SDKConfig ycsdk.Config

	// API address of any YC service
	APIAddress string
}

func GetConnection(ctx context.Context, conf ConnectionConfig, customOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var dialOpts []grpc.DialOption
	sdk, err := ycsdk.Build(ctx, conf.SDKConfig, customOpts...)
	if err != nil {
		return nil, fmt.Errorf("ycsdk.Build: %w", err)
	}
	tokenMiddleware := ycsdk.NewIAMTokenMiddleware(sdk, time.Now)
	dialOpts = append(dialOpts,
		grpc.WithContextDialer(dial.NewProxyDialer(dial.NewDialer())),
		grpc.WithChainUnaryInterceptor(requestid.Interceptor()),
		grpc.WithChainUnaryInterceptor(tokenMiddleware.InterceptUnary),
		grpc.WithChainStreamInterceptor(tokenMiddleware.InterceptStream),
	)

	if conf.SDKConfig.DialContextTimeout > 0 {
		dialOpts = append(dialOpts, grpc.WithBlock(), grpc.WithTimeout(conf.SDKConfig.DialContextTimeout)) // nolint
	}
	if conf.SDKConfig.Plaintext {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := conf.SDKConfig.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		creds := credentials.NewTLS(tlsConfig)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}
	// Append custom options after default, to allow to customize dialer and etc.
	dialOpts = append(dialOpts, customOpts...)

	cc := grpcclient.NewLazyConnContext(grpcclient.DialOptions(dialOpts...))
	conn, err := cc.GetConn(ctx, conf.APIAddress)
	if err != nil {
		return nil, fmt.Errorf("get lazy connection: %w", err)
	}
	return conn, nil
}
