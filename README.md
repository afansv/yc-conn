# yc-conn

[![Go Reference](https://pkg.go.dev/badge/github.com/afansv/yc-conn.svg)](https://pkg.go.dev/github.com/afansv/yc-conn)

Yandex.Cloud Go SDK wrapper for getting gRPC client connection with original SDK auth to use with any YC API (even if
it's not listed in [endpoints](https://yandex.cloud/ru/docs/api-design-guide/concepts/endpoints))

## Why?
Some Yandex.Cloud APIs are not listed in [endpoints](https://yandex.cloud/ru/docs/api-design-guide/concepts/endpoints) for some reason. 
That's why we can't just patch [yc-go-sdk](https://github.com/yandex-cloud/go-sdk) to add new services support

## Installation

```shell
go get -u github.com/afansv/yc-conn
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	ycconn "github.com/afansv/yc-conn"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/speechsense/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

const token = "<YOUR_OAUTH_TOKEN>"

const talkAnalyticsAPIAddr = "api.talk-analytics.yandexcloud.net:443"

func main() {
	ctx := context.Background()

	conn, err := ycconn.GetConnection(ctx, ycconn.ConnectionConfig{
		SDKConfig: ycsdk.Config{
			Credentials: ycsdk.OAuthToken(token),
		},
		APIAddress: talkAnalyticsAPIAddr,
	})
	if err != nil {
		panic(err)
	}

	ss := speechsense.NewTalkServiceClient(conn)

	// Now you can use service client from yandex-cloud/go-genproto
}

```
