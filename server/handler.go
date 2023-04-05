package server

// @Author: Derek
// @Description:
// @Date: 2023/3/22 23:15
// @Version 1.0

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}

type handlerFunc func(c *Context)
