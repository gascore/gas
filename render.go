package gas

// RenderCore render station
type RenderCore struct {
	BE BackEnd
}

// RenderNode node storing changes
type RenderNode struct {
	index int // The index of the item in the heap.

	Type RenderType

	New, Old                     interface{} // *Component, string, int, etc
	NodeParent, NodeNew, NodeOld interface{} // *dom.Element, etc

	Data map[string]interface{} // using only for Type == DataType
}

// RenderType renderNode type
type RenderType int

const (
	// ReplaceType type for replace node
	ReplaceType RenderType = iota

	// CreateType type for create nodes
	CreateType

	// DeleteType type for delete node
	DeleteType

	// RecreateType type for ReCreate
	RecreateType
)

// Add push render nodes to render queue and trying to execute all queue
func (rc *RenderCore) Add(node *RenderNode) {
	err := rc.BE.ExecNode(node)
	if err != nil {
		rc.BE.ConsoleError(err.Error())
	}
}
