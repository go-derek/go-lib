package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// @Author: Derek
// @Description:
// @Date: 2023/3/22 23:14
// @Version 1.0

func TestHandlerBasedOnTree_Route(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	// 要确认已经为支持的方法创建了节点
	assert.Equal(t, len(supportMethods), len(handler.forest))

	postNode := handler.forest[http.MethodPost]

	err := handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(postNode.children))

	n := postNode.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.pattern)
	assert.NotNil(t, n.handler)
	assert.Empty(t, n.children)

	// 我们只有
	//  user -> profile
	err = handler.Route(http.MethodPost, "/user/profile", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(n.children))
	profileNode := n.children[0]
	assert.NotNil(t, profileNode)
	assert.Equal(t, "profile", profileNode.pattern)
	assert.NotNil(t, profileNode.handler)
	assert.Empty(t, profileNode.children)

	// 试试重复
	err = handler.Route(http.MethodPost, "/user", func(c *Context) {})
	assert.Nil(t, err)
	n = postNode.children[0]
	assert.NotNil(t, n)
	assert.Equal(t, "user", n.pattern)
	assert.NotNil(t, n.handler)
	// 有profile节点
	assert.Equal(t, 1, len(n.children))

	// 给 user 再加一个节点
	err = handler.Route(http.MethodPost, "/user/home", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(n.children))
	homeNode := n.children[1]
	assert.NotNil(t, homeNode)
	assert.Equal(t, "home", homeNode.pattern)
	assert.NotNil(t, homeNode.handler)
	assert.Empty(t, homeNode.children)

	// 添加 /order/detail
	err = handler.Route(http.MethodPost, "/order/detail", func(c *Context) {})
	assert.Equal(t, 2, len(postNode.children))
	orderNode := postNode.children[1]
	assert.NotNil(t, orderNode)
	assert.Equal(t, "order", orderNode.pattern)
	// 此刻我们只有/order/detail，但是没有/order
	assert.Nil(t, orderNode.handler)
	assert.Equal(t, 1, len(orderNode.children))

	orderDetailNode := orderNode.children[0]
	assert.NotNil(t, orderDetailNode)
	assert.Empty(t, orderDetailNode.children)
	assert.Equal(t, "detail", orderDetailNode.pattern)
	assert.NotNil(t, orderDetailNode.handler)

	// 加一个 /order
	err = handler.Route(http.MethodPost, "/order", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(postNode.children))
	orderNode = postNode.children[1]
	assert.Equal(t, "order", orderNode.pattern)
	// 此时我们有了 /order
	assert.NotNil(t, orderNode.handler)

	err = handler.Route(http.MethodPost, "/order/*", func(c *Context) {})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(orderNode.children))
	orderWildcard := orderNode.children[1]
	assert.NotNil(t, orderWildcard)
	assert.NotNil(t, orderWildcard.handler)
	assert.Equal(t, "*", orderWildcard.pattern)

	err = handler.Route(http.MethodPost, "/order/*/checkout", func(c *Context) {})
	assert.Equal(t, ErrorInvalidRouterPattern, err)

	err = handler.Route(http.MethodConnect, "/order/checkout", func(c *Context) {})
	assert.Equal(t, ErrorInvalidMethod, err)

	err = handler.Route(http.MethodPost, "/order/:id", func(c *Context) {})
	assert.Nil(t, err)
	// 这时候我们有/order/* 和 /order/:id
	// 因为我们并没有认为它们不兼容，而是/order/:id优先
	assert.Equal(t, 3, len(orderNode.children))
	orderParamNode := orderNode.children[2]
	assert.Equal(t, ":id", orderParamNode.pattern)
}
