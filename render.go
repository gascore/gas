package gas

// RenderCore render station
type RenderCore struct {
	BE BackEnd

	queue []*RenderTask
}

// RenderTask node storing changes
type RenderTask struct {
	Type RenderType

	Parent *Element

	New, Old                     interface{} // *Element, string, int, ...
	NodeParent, NodeNew, NodeOld interface{} // *dom.Element

	ReplaceCanGoDeeper bool
	IgnoreHooks        bool // don't exec elements hooks
}

// RenderType RenderTask type
type RenderType int

const (
	// ReplaceType type for replace node
	ReplaceType RenderType = iota

	// ReplaceHooks type for run after replace hooks
	ReplaceHooks

	// CreateType type for create nodes
	CreateType

	// DeleteType type for delete node
	DeleteType

	// RecreateType type for ReCreate
	RecreateType
)

// Add push render nodes to render queue and trying to execute all queue
func (rc *RenderCore) Add(node *RenderTask) {
	rc.queue = append(rc.queue, node)
}

// GetAll return render nodes from queue
func (rc *RenderCore) GetAll() []*RenderTask {
	return rc.queue
}

// Exec run all render nodes in render core
func (rc *RenderCore) Exec() {
	rc.BE.ExecTasks(rc.queue)
	rc.queue = []*RenderTask{}
}
