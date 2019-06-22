package gas

import (
	"sync"
	"github.com/eapache/queue"
)

// RenderCore render station
type RenderCore struct {
	Queue *queue.Queue
	BE    BackEnd

	WG *sync.WaitGroup
}

// RenderNode node storing changes
type RenderNode struct {
	index int // The index of the item in the heap.

	Type     RenderType

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
func (rc *RenderCore) Add(nodes []*RenderNode) {
	rc.WG.Wait()
	rc.WG.Add(1)
	go func() {
		for _, node := range nodes {
			rc.Queue.Add(node)
		}
		rc.WG.Done()
	}()

	// trying to execute all renderNodes in queue
	rc.WG.Wait()
	for rc.Queue.Length() > 0 {
		err := rc.BE.ExecNode(rc.Queue.Remove().(*RenderNode))
		if err != nil {
			rc.BE.ConsoleError(err.Error())
			return
		}
	}
}

func singleNode(node *RenderNode) []*RenderNode {
	return []*RenderNode{node}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes	[]*RenderNode
	head	int
	tail	int
	count	int
}
