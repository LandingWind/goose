package goose

type TrieNode struct {
	fullPath   string      // 完整path
	curPath    string      // 当前节点对应部分的path
	childNodes []*TrieNode // 当前节点的子节点
	isDynamic  bool        // 是否为动态路由节点 即如/:name
}

func (node *TrieNode) insertPath(fullPath string, parts []string, height int) {
	if len(parts) == height {
		node.fullPath = fullPath
		return
	}
	curPart := parts[height] // 当前短路径
	// 寻找当前Trie树是否有该短路径
	child := node.matchChild(curPart)
	if child == nil {
		// 没有该短路径 则创建新短路径
		isDynamic := false
		if curPart[0] == ':' {
			isDynamic = true
		}
		child = &TrieNode{
			curPath:    curPart,
			childNodes: nil,
			isDynamic:  isDynamic,
			fullPath:   "",
		}
		node.childNodes = append(node.childNodes, child)
	}
	// 下一层短路径
	child.insertPath(fullPath, parts, height+1)
}

func (node *TrieNode) matchChild(curPart string) *TrieNode {
	for _, val := range node.childNodes {
		if val.curPath == curPart {
			return val
		}
	}
	return nil
}

/**
* 查找是否有对应路径
* 需要满足可以在trie树内完整匹配pattern，而且最终node含有fullPath
 */
func (node *TrieNode) searchPath(parts []string, height int) *TrieNode {
	if len(parts) == height {
		if node.fullPath == "" {
			return nil
		}
		return node
	}
	curPart := parts[height]
	children := node.matchChildren(curPart)
	for _, child := range children {
		result := child.searchPath(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
func (node *TrieNode) matchChildren(curPart string) []*TrieNode {
	res := make([]*TrieNode, 0)
	for _, val := range node.childNodes {
		if val.curPath == curPart || val.isDynamic {
			res = append(res, val)
		}
	}
	return res
}
