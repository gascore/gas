package gas

import (
	"sync"
)

// RenderCore render station
type RenderCore struct {
	Queue *PriorityQueue
	BE    BackEnd

	WG *sync.WaitGroup
}

// RenderNode node storing changes
type RenderNode struct {
	index int // The index of the item in the heap.

	Type     RenderType
	Priority Priority

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

	// DataType type for Set, SetValue
	DataType

	// SyncType type for update g-model value, remove g-show styles
	SyncType

	// RecreateType type for ReCreate
	RecreateType
)

// Priority RenderNode priority (the more the more important)
type Priority int

const (
	// EventPriority pritority for Set, SetValue
	EventPriority Priority = iota

	// InputPriority priority for g-model input events
	InputPriority

	// RenderPriority priority for Create, Replace, Delete, ForceUpdate, ReCreate
	RenderPriority
)

// Add push render nodes to render queue and trying to execute all queue
func (rc *RenderCore) Add(nodes []*RenderNode) {
	rc.WG.Wait()
	rc.WG.Add(1)
	go func() {
		for _, node := range nodes {
			rc.Queue.Push(node)
		}
		rc.WG.Done()
	}()

	// trying to execute all renderNodes in queue
	rc.WG.Wait()
	for rc.Queue.Len() > 0 {
		node := rc.Queue.Pop().(*RenderNode)

		switch node.Type {
		case DataType:
			newC, ok := node.New.(*Component)
			if !ok {
				rc.BE.ConsoleError("invalid New type in RenderNode with DataType")
				return
			}

			err := newC.realSet(node)
			if err != nil {
				rc.BE.ConsoleError(err.Error())
				return
			}
			break
		default:
			err := rc.BE.ExecNode(node)
			if err != nil {
				rc.BE.ConsoleError(err.Error())
				return
			}
		}
	}
}

func singleNode(node *RenderNode) []*RenderNode {
	return []*RenderNode{node}
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*RenderNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

// Swap swap two nodes in queue
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push push node to queue
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*RenderNode)
	item.index = n
	*pq = append(*pq, item)
}

// Pop pop node from queue
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	item := old[0]
	item.index = -1 // for safety
	*pq = old[1:]
	return item
}
