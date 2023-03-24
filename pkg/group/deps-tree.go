package group

import (
	"fmt"
	"strings"
)

// DepsGraph is (intended to be) a directed acyclic graph that tracks the
// requirements dependencies between phases within a goal.
type DepsGraph struct {
	nodes map[string]struct{}
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
// index into the slice of the edge at from. Returns -1 if no such edge exists.
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

// NewDepsGraph constructs a DepsGraph or the given goal. Each task names the
// phases that must run before the phase to which the task belongs. The
// constructed directed acyclic graph (DAG) has edges pointing from the phase
// that must run earlier to the phases that may run only after that phase runs.
func NewDepsGraph(goal *Goal) *DepsGraph {
	nodes := make(map[string]struct{}, len(goal.Phases))
	edges := make(map[string][]string, len(goal.Phases))

	for _, phase := range goal.Phases {
		nodes[phase.Name] = struct{}{}

		for _, task := range phase.Tasks() {
			for _, req := range task.Requires() {
				addEdge(edges, phase.Name, req)
			}
		}
	}

	return &DepsGraph{nodes, edges}
}

// PhaseOrder constructs an ordered list of phases based upon the DAG stored in
// DepsGraph. In the DAG, the most required nodes (vertices) will have no edges
// directed to them. Therefore, to produce a set of tasks that are grouped and
// ordered, we do the following:
//
//  1. Find all unmarked nodes that have no edges pointing to them by unmarked
//     nodes. This becomes the set of phases that may run. These will be ordered by
//     name if there is more than one.
//
//  2. Mark all the nodes found in (1).
//
//  3. Repeat (1) and (2) until all nodes are marked.
//
// The result is a slice of phase names in the order the phases must be run.
//
// Along the way, it is possible for a cycle to be detected in the graph. If
// this occurs, it means that one or more task descriptions is misconfigured or
// that the mix of plugins have incompatible dependency requirements. In such a
// case, the error is reported with a description of a cycle.
func (d *DepsGraph) PhaseOrder() ([]string, error) {
	workEdges := make(map[string][]string, len(d.edges))
	copyEdges(workEdges, d.edges)

	workNodes := make(map[string]struct{}, len(d.nodes))
	for k := range d.nodes {
		workNodes[k] = struct{}{}
	}

	out := make([]string, 0, len(d.nodes))
	foundNodes := make([]string, 0, len(d.nodes))
	for len(workNodes) > 0 {
		for from := range workNodes {
			if _, hasAnyEdges := workEdges[from]; hasAnyEdges {
				continue
			}

			foundNodes = append(foundNodes, from)
		}

		// this would mean there's a cycle in the dependency graph
		if len(foundNodes) == 0 {
			// remaining nodes form the cycle, so let's write them out
			cycle := &strings.Builder{}
			for from := range workNodes {
				for _, to := range workEdges[from] {
					_, _ = fmt.Fprintf(cycle, " - %s => %s\n", from, to)
				}
			}
			return nil, fmt.Errorf("tasks contain a dependency cycle, cannot run this goal:\n%s", cycle.String())
		}

		// once used, we can remove these from the list
		for _, found := range foundNodes {
			out = append(out, found)
			delete(workNodes, found)
			for from := range workNodes {
				deleteEdge(workEdges, from, found)
			}
		}

		// reuse the previous memory
		foundNodes = foundNodes[:0]
	}

	return out, nil
}
