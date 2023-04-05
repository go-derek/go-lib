package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// @Author: Derek
// @Description:
// @Date: 2023/3/19 10:12
// @Version 1.0

type Server interface {
	Routable
	Start(address string) (err error)
	Shutdown(ctx context.Context) error
}

// sdkHttpServer 这个是基于 net/http 这个包实现的 http server
type sdkHttpServer struct {
	// Name server 的名字，给个标记，日志输出的时候用得上
	Name    string
	handler Handler
	root    Filter
	ctxPool sync.Pool
}

func (s *sdkHttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := s.ctxPool.Get().(*Context)
	defer func() {
		s.ctxPool.Put(c)
	}()
	c.Reset(writer, request)
	s.root(c)
}

func (s *sdkHttpServer) Route(method string, pattern string, handlerFunc handlerFunc) error {
	return s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) (err error) {
	return http.ListenAndServe(address, s)
}

func (s *sdkHttpServer) Shutdown(ctx context.Context) error {
	fmt.Printf("%s shutdown start\n", s.Name)
	time.Sleep(time.Second)
	fmt.Printf("%s shutdown end\n", s.Name)
	return nil
}

func NewSdkHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBasedOnTree()

	// 因为我们是一个链，所以我们把最后的业务逻辑处理，也作为一环
	var root Filter = handler.ServeHTTP

	// 从后往前把filter串起来
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	res := &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
		ctxPool: sync.Pool{New: func() interface{} {
			return newContext()
		}},
	}
	return res
}
