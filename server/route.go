package server

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

// @Author: Derek
// @Description:
// @Date: 2023/3/19 10:13
// @Version 1.0

var (
	ErrorInvalidRouterPattern = errors.New("invalid router pattern")
	ErrorInvalidMethod        = errors.New("invalid method")
)

// Routable 可路由的
type Routable interface {
	// Route 设定一个路由，命中该路由的会执行handlerFunc的代码
	Route(method string, pattern string, handlerFunc handlerFunc) error
}

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

type HandlerBasedOnTree struct {
	forest map[string]*node
}

func NewHandlerBasedOnTree() Handler {
	forest := make(map[string]*node, len(supportMethods))
	for _, m := range supportMethods {
		forest[m] = newRootNode(m)
	}
	return &HandlerBasedOnTree{
		forest: forest,
	}
}

// ServeHTTP 就是从树里面找节点
// 找到了就执行
func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	handler, found := h.findRouter(c.R.Method, c.R.URL.Path, c)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)
}

func (h *HandlerBasedOnTree) Route(method string, pattern string, handlerFunc handlerFunc) (err error) {
	err = h.validatePattern(pattern)
	if err != nil {
		return err
	}

	// 将pattern按照URL的分隔符切割
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	// 当前指向根节点
	cur, ok := h.forest[method]
	if !ok {
		return ErrorInvalidMethod
	}

	for index, path := range paths {
		// 匹配path节点
		matchChild, found := h.findMatchChild(cur, path, nil)
		// != nodeTypeAny 是考虑到 /order/* 和 /order/:id 这种注册顺序
		if found && matchChild.ty != nodeTypeAny {
			cur = matchChild
		} else {
			// 为当前节点根据
			h.createSubTree(cur, paths[index:], handlerFunc)
			return
		}
	}
	// 离开了循环，说明我们加入的是短路径，
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handlerFunc

	return
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handlerFn handlerFunc) {
	cur := root
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFn
}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string, c *Context) (*node, bool) {
	candidates := make([]*node, 0, 2)
	for _, child := range root.children {
		if child.matchFunc(path, c) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].ty < candidates[j].ty
	})
	return candidates[len(candidates)-1], true
}

func (h *HandlerBasedOnTree) validatePattern(pattern string) (err error) {
	pos := strings.Index(pattern, "*")
	if pos > 0 {
		if pos != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return
}

func (h *HandlerBasedOnTree) findRouter(method, path string, c *Context) (handlerFunc, bool) {
	cur, ok := h.forest[method]
	if !ok {
		return nil, false
	}

	paths := strings.Split(strings.Trim(path, "/"), "/")

	for _, p := range paths {
		matchChild, found := h.findMatchChild(cur, p, c)
		if !found {
			return nil, false
		}
		cur = matchChild
	}

	if cur.handler == nil {
		return nil, false
	}

	return cur.handler, true
}
