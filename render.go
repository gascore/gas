package gas

// RenderCore render station
type RenderCore struct {
	BE BackEnd

	queue   []*RenderTask
	counter int64
}

// RenderTask node storing changes
type RenderTask struct {
	ID int64

	Type RenderType

	Parent *Element

	New, Old                     interface{} // *Element, string, int, ...
	NodeParent, NodeNew, NodeOld interface{} // *dom.Element

	ReplaceCanGoDeeper bool
	InReplaced         bool
}

// RenderType RenderTask type
type RenderType int

const (
	// RReplace type for replace node
	RReplace RenderType = iota

	// RReplaceHooks type for run after replace hooks
	RReplaceHooks

	// RCreate type for create nodes
	RCreate

	// RFirstRender type for first gas render
	RFirstRender

	// RDelete type for delete node
	RDelete

	// RRecreate type for recreate node
	RRecreate
)

// Add push render tasks to render queue and trying to execute all queue
func (rc *RenderCore) Add(task *RenderTask) {
	task.ID = rc.counter
	rc.queue = append(rc.queue, task)
	rc.counter++
}

// GetAll return render nodes from queue
func (rc *RenderCore) GetAll() []*RenderTask {
	return rc.queue
}

// Exec run all render nodes in render core
func (rc *RenderCore) Exec() {
	tQueue := rc.queue
	rc.queue = []*RenderTask{}
	rc.BE.ExecTasks(tQueue)
}
