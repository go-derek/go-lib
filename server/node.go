package server

import "strings"

// @Author: Derek
// @Description:
// @Date: 2023/3/22 23:18
// @Version 1.0

const patternAny = "*"

// nodeType
const (

	// 根节点，只有根用这个
	nodeTypeRoot nodeType = iota

	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// 正则
	nodeTypeReg

	// 静态，即完全匹配
	nodeTypeStatic
)

type (
	node struct {
		children []*node

		// 如果这是叶子节点，那么匹配上之后就可以调用 handler
		handler   handlerFunc
		matchFunc matchFunc

		// 匹配到当前节点的 pattern
		pattern string
		// ty 类型&&决定优先级
		ty nodeType
	}
	// matchFunc 承担两个职责，一个是判断是否匹配，一个是在匹配之后
	// 将必要的数据写入到 Context
	// 所谓必要的数据，这里基本上是指路径参数
	matchFunc func(path string, c *Context) bool
	nodeType  int
)

func newRootNode(method string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			panic("never call me")
		},
		ty:      nodeTypeRoot,
		pattern: method,
	}
}

// return new node
func newNode(path string) *node {
	if path == "*" {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

// 通配符 * 节点
func newAnyNode() *node {
	return &node{
		// 因为我们不允许 * 后面还有节点，所以这里可以不用初始化
		//children: make([]*node, 0, 2),
		matchFunc: func(string, *Context) bool {
			return true
		},
		ty:      nodeTypeAny,
		pattern: patternAny,
	}
}

// 路径参数节点
func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			if c != nil {
				c.PathParams[paramName] = p
			}
			// 如果自身是一个参数路由，
			// 然后又来一个通配符，我们认为是不匹配的
			return p != patternAny
		},
		ty:      nodeTypeParam,
		pattern: path,
	}
}

// 静态节点
func newStaticNode(path string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, _ *Context) bool {
			return path == p && p != "*"
		},
		ty:      nodeTypeStatic,
		pattern: path,
	}
}
