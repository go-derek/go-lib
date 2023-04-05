package server

import (
	"context"
	"net/http"
)

// @Author: Derek
// @Description:
// @Date: 2023/3/19 10:24
// @Version 1.0

type httpServer struct {
}

func (s *httpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *httpServer) Route(method string, pattern string, handlerFunc handlerFunc) error {
	//TODO implement me
	panic("implement me")
}

func (s *httpServer) Start(addr string) (err error) {
	return http.ListenAndServe(addr, s)
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewHttpServer() Server {

	return &httpServer{}
}
