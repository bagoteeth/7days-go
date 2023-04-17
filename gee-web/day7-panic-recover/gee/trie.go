package gee

import "strings"

type node struct {
	pattern string
	//树的当前路由部分
	part     string
	children []*node
	//不精确匹配时为true
	isWild bool
}

func (r *node) matchChild(part string) *node {
	for _, c := range r.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

func (r *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, c := range r.children {
		if c.part == part || c.isWild {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

//parts是一个路由的所有层按顺序排列，一般为pattern split by "/"
//height是当前匹配的层级，从0开始不断自增1，代表r的层级，最终len(parts) == height时才是最终匹配的
func (r *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		r.pattern = pattern
		return
	}
	part := parts[height]
	child := r.matchChild(part)
	if child == nil {
		child = &node{
			pattern:  "",
			part:     part,
			children: nil,
			isWild:   part[0] == ':' || part[0] == '*',
		}
		r.children = append(r.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//完整匹配(/a/b/c <=> /a/b/c)或匹配到通配动态路由(/a/b/c <=> /a/*path)
func (r *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(r.part, "*") {
		if r.pattern == "" {
			return nil
		}
		return r
	}
	part := parts[height]
	children := r.matchChildren(part)
	for _, c := range children {
		res := c.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
