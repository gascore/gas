package gas

import (
	"sync"
)

type RenderCore struct {
	Queue []*RenderNode
	BE    BackEnd

	WG *sync.WaitGroup
	M  *sync.Mutex
}

type RenderNode struct {
	Type RenderType

	New, Old interface{} // *Component, string, int, etc
	NodeParent, NodeNew, NodeOld interface{} // *dom.Element, etc
}

type RenderType int
const (
	ReplaceType RenderType = iota
	CreateType
	DeleteType
	SyncType // update g-model value, remove g-show styles
	RecreateType // ReCreate
)

func (rc *RenderCore) Add(node *RenderNode) {
	rc.WG.Add(1)
	go func() {
		rc.M.Lock()

		rc.Queue = append(rc.Queue, node)

		rc.M.Unlock()
		rc.WG.Done()
	}()
}

func (rc *RenderCore) AddMany(nodes []*RenderNode) {
	for _, node := range nodes {
		rc.Add(node)
	}
}

func (rc *RenderCore) Run() {
	go func() {
		rc.WG.Wait()
		rc.M.Lock()

		for len(rc.Queue) != 0 {
			node := rc.Queue[0]

			err := rc.BE.ExecNode(node)
			if err != nil {
				rc.M.Unlock()
				rc.BE.ConsoleError(err.Error())
				return
			}

			// remove first element from queue
			copy(rc.Queue[0:], rc.Queue[1:])
			rc.Queue[len(rc.Queue)-1] = nil
			rc.Queue = rc.Queue[:len(rc.Queue)-1]
		}

		rc.M.Unlock()
	}()
}

func aloneNode(node *RenderNode) []*RenderNode {
	return []*RenderNode{node}
}