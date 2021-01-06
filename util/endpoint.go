package util

import (
	"context"
	"net/http"
)

//请求
type Endpoint func(ctx context.Context, requestParam interface{}) (response interface{}, err error)
//决定请求路径
type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error
