package master

import (
	"fmt"
	"strings"

	"github.com/zostay/zedpm/plugin"
)

// DepNode is a vertex in the directed acyclic graph that is used to track task
// requirements for use in ordering tasks.
type DepNode struct {
	// Path is the task path for a task.
	Path string

	// Tasks is the list of tasks that the task named by Path requires to run
	// before it runs.
	Tasks []plugin.TaskDescription
}

// DepsGraph is (intended to be) a directed acyclic graph that tracks the
// requirements dependencies between tasks.
type DepsGraph struct {
	root  string
	nodes map[string]*DepNode
	edges map[string][]string
}

// addEdge adds a new edge to the DAG.
func addEdge(edges map[string][]string, from, to string) {
	if _, hasEdgeFrom := edges[from]; hasEdgeFrom {
		for _, e := range edges[from] {
			if e == to {
				return
			}
		}
		edges[from] = append(edges[from], to)
	} else {
		edges[from] = []string{to}
	}
}

// edgeIndex finds the index of the given edge in the DAG--where index means the
// index into the slice of the the edge at from. Returns -1 if no such edge
// exists.
func edgeIndex(edges map[string][]string, from, to string) int {
	if _, hasEdgeFrom := edges[from]; !hasEdgeFrom {
		return -1
	}

	for i, toCheck := range edges[from] {
		if toCheck == to {
			return i
		}
	}

	return -1

}

// deleteEdge deletes an edge from the DAG.
func deleteEdge(edges map[string][]string, from, to string) {
	for {
		i := edgeIndex(edges, from, to)
		if i < 0 {
			return
		}

		if len(edges[from]) == 1 {
			delete(edges, from)
			return
		}

		copy(edges[from][i:], edges[from][i+1:])
		edges[from] = edges[from][:len(edges[from])-1]
	}
}

// copyEdges creates a deep copy of the DAG edges.
func copyEdges(dst, src map[string][]string) {
	for from, tos := range src {
		dst[from] = make([]string, len(tos))
		copy(dst[from], tos)
	}
}

// NewDepsGraph constructs a DepsGraph for the named goal with the given tasks.
// This immediately configures the direct acyclic graph (DAG) kept by the
// DepsGraph object.
//
// This DAG is setup with two sets of dependencies:
//
//  1. Hierarchical dependencies exist to goal from task from sub-task from
//     sub-sub-task and so on. That is, a sub-sub-task has a directed edge to its
//     parent sub-task, which has a directed edge to its parent task, which has a
//     directed edge to its parent goal.
//
//  2. Requirements dependencies exist from required task to task doing the
//     requiring. That is, if /release/publish requires /release/mint, there
//     exists an edge in the DAG pointing fro /release/mint to /release/publish.
//
// The contained DAG has the relationship of requirements inverted. The most
// required node will have no edges pointing to it. Anything that requires
// something else will have one or more edges pointing at it in the graph.
func NewDepsGraph(goal string, tasks []plugin.TaskDescription) *DepsGraph {
	nodes := make(map[string]*DepNode, len(tasks))
	edges := make(map[string][]string, len(tasks))
	goalPath := "/" + goal

	nodes[goalPath] = &DepNode{goalPath, []plugin.TaskDescription{}}

	establishNode := func(tree string) {
		result := nodes[tree]
		if result == nil {
			result = &DepNode{
				Path:  tree,
				Tasks: make([]plugin.TaskDescription, 0, 1),
			}
			nodes[tree] = result
		}
	}

	for _, task := range tasks {
		parts := strings.Split(task.Name()[1:], "/")

		// skip [0], that's the goal path
		taskNames := parts[1:]

		parent := goalPath
		tree := parent
		for _, part := range taskNames {
			tree += "/" + part

			establishNode(tree)

			addEdge(edges, parent, tree)

			parent = tree
		}

		nodes[task.Name()].Tasks = append(nodes[task.Name()].Tasks, task)

		for _, req := range task.Requires() {
			addEdge(edges, task.Name(), req)
		}
	}

	return &DepsGraph{goalPath, nodes, edges}
}

// GroupOrder constructs an ordered list of plugin.TaskDescription groups based
// upon the DAG stored in DepsGraph. In the DAG, the most required ndoes
// (verticed) will have no edges directed to them. Therefore, to produce a set
// of tasks that are grouped and ordered, we do the following:
//
//  1. Find all unmarked nodes that have no edges pointing to them by unmarked
//     nodes. This becomes the next slice of plugin.TaskDescription objects.
//
//  2. Mark all the nodes found in (1).
//
//  3. Repeat (1) and (2) until all nodes are marked.
//
// The result is a slice of slices of plugin.TaskDescription objects, which
// represent tasks that can safely be run concurrently.
//
// Along the way, at least two error conditions could be encountered, which sill
// result in this method failing with an error:
//
//  1. If an edge is encountered that refers to a node that does not exist (e.g.,
//     a plugin.TaskDescription has a requirement that belongs to another goal),
//     an error is returned. A failed requirement indicates either a issing
//     plugin dependency or a plugin that contains a serious bug.
//
//  2. If while looking for unmarked nodes that have no edges poining to them by
//     unmarked nodes we run into a case where we get zero such nodes, but the
//     number of unmarked nodes is non-zero, then we have a cycle. This is a
//     directed acyclic graph. A cycle indicates that some plugin contains a
//     serious bug.
func (d *DepsGraph) GroupOrder() ([][]plugin.TaskDescription, error) {
	workEdges := make(map[string][]string, len(d.edges))
	copyEdges(workEdges, d.edges)

	workNodes := make(map[string]*DepNode, len(d.nodes))
	for k, v := range d.nodes {
		workNodes[k] = v
	}

	out := make([][]plugin.TaskDescription, 0, len(d.nodes))
	foundNodes := make([]string, 0, len(d.nodes))
	for len(workNodes) > 0 {
		for from := range workNodes {
			if edges, hasAnyEdges := workEdges[from]; hasAnyEdges {
				// detect nodes that get mentioned, but aren't defined
				for _, edge := range edges {
					if _, edgePointsToRealNode := workNodes[edge]; !edgePointsToRealNode {
						return nil, fmt.Errorf("tasks contain an unfulfilled dependency %q, cannot run this goal", edge)
					}
				}
				continue
			}

			// this is what we're really searching for
			foundNodes = append(foundNodes, from)
		}

		// this would mean there's a cycle in the dependency graph
		if len(foundNodes) == 0 {
			// remaining nodes form the cycle, so let's write them out
			cycle := &strings.Builder{}
			for from := range workNodes {
				for _, to := range workEdges[from] {
					fmt.Fprintf(cycle, " - %s => %s\n", from, to)
				}
			}
			return nil, fmt.Errorf("tasks contain a dependency cycle, cannot run this goal:\n%s", cycle.String())
		}

		// once used, we can remove these from the list
		thisPhase := make([]plugin.TaskDescription, 0, len(foundNodes))
		for _, used := range foundNodes {
			thisPhase = append(thisPhase, workNodes[used].Tasks...)
			delete(workNodes, used)
			for from := range workNodes {
				deleteEdge(workEdges, from, used)
			}
		}

		// some phases may have no tasks to execute, but we don't care about
		// those
		if len(thisPhase) > 0 {
			out = append(out, thisPhase)
		}

		// reuse the previous memory
		foundNodes = foundNodes[:0]
	}

	return out, nil
}
