package gas

import (
	"sync"
)

// RenderCore
type RenderCore struct {
	Queue *PriorityQueue
	BE    BackEnd

	WG *sync.WaitGroup
}

// RenderNode
type RenderNode struct {
	index int // The index of the item in the heap.

	Type RenderType
	Priority Priority

	New, Old interface{} // *Component, string, int, etc
	NodeParent, NodeNew, NodeOld interface{} // *dom.Element, etc

	Data map[string]interface{} // using only for Type == DataType
}

// RenderType renderNode type
type RenderType int
const (
	ReplaceType RenderType = iota
	CreateType
	DeleteType
	DataType // Set, SetValue
	SyncType // update g-model value, remove g-show styles
	RecreateType // ReCreate
)

// RenderType renderNode priority (the more the more important)
type Priority int
const (
	EventPriority Priority = iota // Set, SetValue
	RenderPriority // Create, Replace, Delete, ForceUpdate, ReCreate
	InputPriority // Using in g-model input events
)

func (rc *RenderCore) Add(nodes []*RenderNode) {
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
				// rc.M.Unlock()
				rc.BE.ConsoleError("invalid New type in RenderNode with DataType")
				return
			}

			err := newC.realSet(node)
			if err != nil {
				// rc.M.Unlock()
				rc.BE.ConsoleError(err.Error())
				return
			}
		default:
			err := rc.BE.ExecNode(node)
			if err != nil {
				// rc.M.Unlock()
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

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*RenderNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	item := old[0]
	item.index = -1 // for safety
	*pq = old[1:]
	return item
}