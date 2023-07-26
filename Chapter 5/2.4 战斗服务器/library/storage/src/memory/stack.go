package memory

type stack struct {
	nodes []*exprNode
}

func newStack() *stack {
	return &stack{nodes: []*exprNode{}}
}

func (stack *stack) Push(value ...*exprNode) {
	stack.nodes = append(stack.nodes, value...)
}

func (stack *stack) Top() *exprNode {
	length := len(stack.nodes)
	if length > 0 {
		return stack.nodes[length-1]
	}
	return nil
}

func (stack *stack) Pop() *exprNode {
	length := len(stack.nodes)
	if length > 0 {
		last := stack.nodes[length-1]
		stack.nodes = stack.nodes[:length-1]
		return last
	}
	return nil
}

func (stack *stack) Size() int {
	return len(stack.nodes)
}
