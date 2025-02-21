package prefix_tree

type Node struct {
	children    map[rune]*Node
	isEndOfWord bool
}

func NewNode() *Node {
	return &Node{
		children:    make(map[rune]*Node),
		isEndOfWord: false,
	}
}

type PrefixTree struct {
	root *Node
}

func NewTree() *PrefixTree {
	return &PrefixTree{
		root: NewNode(),
	}
}

func (r *PrefixTree) Insert(word string) {
	currentNode := r.root

	for _, c := range word {
		_, exists := currentNode.children[c]
		if !exists {
			currentNode.children[c] = NewNode()
		}
		currentNode = currentNode.children[c]
	}

	currentNode.isEndOfWord = true
}

func (r *PrefixTree) Search(word string) bool {
	currentNode := r.root

	for _, c := range word {
		_, exists := currentNode.children[c]
		if !exists {
			return false
		}
		currentNode = currentNode.children[c]
	}
	return currentNode.isEndOfWord
}

func (r *PrefixTree) StartsWith(word string) bool {
	currentNode := r.root

	for _, c := range word {
		_, exists := currentNode.children[c]
		if !exists {
			return false
		}
		currentNode = currentNode.children[c]
	}
	return true
}

func (r *PrefixTree) GetAllWordsStartingWith(prefix string) []string {
	currentNode := r.root

	for _, c := range prefix {
		_, exists := currentNode.children[c]
		if !exists {
			return []string{}
		}
		currentNode = currentNode.children[c]
	}

	words := []string{}
	if currentNode.isEndOfWord {
		words = append(words, string(prefix))
	}

	suffixes := r.getAllWordsFromNode(currentNode)
	for i := range suffixes {
		suffixes[i] = prefix + suffixes[i]
	}
	words = append(words, suffixes...)

	return words
}

// Tree has apple, applet GetAllWordsStartingWith(app) should return [[apple, applet]]
// Tree has apple, applet, appke, appkel, appde GetAllWordsStartingWith(app) should return [[apple, applet], [appke, appkel], [appde]]
func (r *PrefixTree) GetAllWordsStartingWithGroupedByChildren(prefix string) [][]string {
	currentNode := r.root

	// Navigate to prefix node
	for _, c := range prefix {
		_, exists := currentNode.children[c]
		if !exists {
			return [][]string{}
		}
		currentNode = currentNode.children[c]
	}

	result := [][]string{}
	for c, child := range currentNode.children {
		group := []string{}

		if child.isEndOfWord {
			group = append(group, prefix+string(c))
		}

		suffixes := r.getAllWordsFromNode(child)
		for _, suffix := range suffixes {
			group = append(group, prefix+string(c)+suffix)
		}

		if len(group) > 0 {
			result = append(result, group)
		}
	}

	return result
}

func (r *PrefixTree) getAllWordsFromNode(node *Node) []string {
	words := []string{}

	for c, child := range node.children {
		if child.isEndOfWord {
			words = append(words, string(c))
		}

		suffixes := r.getAllWordsFromNode(child)
		for i := range suffixes {
			suffixes[i] = string(c) + suffixes[i]
		}
		words = append(words, suffixes...)
	}

	return words
}
