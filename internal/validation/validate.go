package validation

import (
	"fmt"

	"github.com/sketch-down/sketch-down/internal/diagram"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Warning struct {
	Severity Severity `json:"severity"`
	NodeID   string   `json:"nodeId,omitempty"`
	EdgeID   string   `json:"edgeId,omitempty"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
}

// Validate runs all checks on a diagram and returns warnings.
func Validate(d *diagram.Diagram) []Warning {
	if d == nil {
		return nil
	}

	var ws []Warning
	nodeMap := buildNodeMap(d)
	portToNode := buildPortToNodeMap(d)

	ws = append(ws, checkMissingVoltage(d)...)
	ws = append(ws, checkMissingFuseRating(d)...)
	ws = append(ws, checkVoltageMismatch(d, portToNode, nodeMap)...)
	ws = append(ws, checkCableNoFuse(d, portToNode, nodeMap)...)
	ws = append(ws, checkCableUndersized(d, portToNode, nodeMap)...)
	ws = append(ws, checkVoltageDrop(d, portToNode, nodeMap)...)
	ws = append(ws, checkFuseOversized(d)...)

	return ws
}

func buildNodeMap(d *diagram.Diagram) map[string]*diagram.Node {
	m := make(map[string]*diagram.Node, len(d.Nodes))
	for i := range d.Nodes {
		m[d.Nodes[i].ID] = &d.Nodes[i]
	}
	return m
}

func buildPortToNodeMap(d *diagram.Diagram) map[string]*diagram.Node {
	m := make(map[string]*diagram.Node)
	for i := range d.Nodes {
		for _, p := range d.Nodes[i].Ports {
			m[p.ID] = &d.Nodes[i]
		}
	}
	return m
}

// hasVoltageField returns true if this node type should have a voltage value.
func hasVoltageField(nt diagram.NodeType) bool {
	switch nt {
	case diagram.NodeBattery, diagram.NodeAlternator, diagram.NodeMPPT, diagram.NodeDCLoad:
		return true
	}
	return false
}

func checkMissingVoltage(d *diagram.Diagram) []Warning {
	var ws []Warning
	for i := range d.Nodes {
		n := &d.Nodes[i]
		if hasVoltageField(n.Type) && diagram.GetNodeVoltage(n) == 0 {
			ws = append(ws, Warning{
				Severity: SeverityWarning,
				NodeID:   n.ID,
				Code:     "missing_voltage",
				Message:  fmt.Sprintf("%s %q has no voltage set", n.Type, n.Label),
			})
		}
	}
	return ws
}

func checkMissingFuseRating(d *diagram.Diagram) []Warning {
	var ws []Warning
	for i := range d.Nodes {
		n := &d.Nodes[i]
		if n.Type == diagram.NodeFuse && diagram.GetNodeMaxAmps(n) == 0 {
			ws = append(ws, Warning{
				Severity: SeverityWarning,
				NodeID:   n.ID,
				Code:     "missing_fuse_rating",
				Message:  fmt.Sprintf("Fuse %q has no rating set", n.Label),
			})
		}
	}
	return ws
}

func checkVoltageMismatch(d *diagram.Diagram, portToNode map[string]*diagram.Node, _ map[string]*diagram.Node) []Warning {
	var ws []Warning
	for _, e := range d.Edges {
		srcNode := portToNode[e.SourcePortID]
		tgtNode := portToNode[e.TargetPortID]
		if srcNode == nil || tgtNode == nil {
			continue
		}

		srcV := diagram.GetNodeVoltage(srcNode)
		tgtV := diagram.GetNodeVoltage(tgtNode)
		if srcV > 0 && tgtV > 0 && srcV != tgtV {
			ws = append(ws, Warning{
				Severity: SeverityError,
				EdgeID:   e.ID,
				Code:     "voltage_mismatch",
				Message:  fmt.Sprintf("Voltage mismatch: %q (%.1fV) connected to %q (%.1fV)", srcNode.Label, srcV, tgtNode.Label, tgtV),
			})
		}
	}
	return ws
}

func checkCableNoFuse(d *diagram.Diagram, portToNode map[string]*diagram.Node, nodeMap map[string]*diagram.Node) []Warning {
	var ws []Warning
	// Build adjacency from ports
	adj := buildAdjacency(d, portToNode)

	for i := range d.Nodes {
		n := &d.Nodes[i]
		if n.Type != diagram.NodeBattery {
			continue
		}
		// BFS from battery, check if we reach a load without passing through a fuse
		visited := map[string]bool{n.ID: true}
		queue := []string{n.ID}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			curNode := nodeMap[cur]
			if curNode == nil {
				continue
			}
			if curNode.Type == diagram.NodeFuse {
				continue // fuse stops traversal on this branch
			}
			for _, neighbor := range adj[cur] {
				if visited[neighbor.nodeID] {
					continue
				}
				visited[neighbor.nodeID] = true
				neighborNode := nodeMap[neighbor.nodeID]
				if neighborNode == nil {
					continue
				}
				if neighborNode.Type == diagram.NodeDCLoad {
					ws = append(ws, Warning{
						Severity: SeverityWarning,
						EdgeID:   neighbor.edgeID,
						Code:     "cable_no_fuse",
						Message:  fmt.Sprintf("Path from battery %q to load %q has no fuse", n.Label, neighborNode.Label),
					})
				}
				queue = append(queue, neighbor.nodeID)
			}
		}
	}
	return ws
}

type adjacencyEntry struct {
	nodeID string
	edgeID string
}

func buildAdjacency(d *diagram.Diagram, portToNode map[string]*diagram.Node) map[string][]adjacencyEntry {
	adj := make(map[string][]adjacencyEntry)
	for _, e := range d.Edges {
		src := portToNode[e.SourcePortID]
		tgt := portToNode[e.TargetPortID]
		if src == nil || tgt == nil {
			continue
		}
		adj[src.ID] = append(adj[src.ID], adjacencyEntry{nodeID: tgt.ID, edgeID: e.ID})
		adj[tgt.ID] = append(adj[tgt.ID], adjacencyEntry{nodeID: src.ID, edgeID: e.ID})
	}
	return adj
}

func checkCableUndersized(d *diagram.Diagram, portToNode map[string]*diagram.Node, _ map[string]*diagram.Node) []Warning {
	var ws []Warning
	for _, e := range d.Edges {
		if e.CableSpec == nil {
			continue
		}
		awg, ok := LookupAWG(e.CableSpec.GaugeAWG)
		if !ok {
			continue
		}

		// Determine current: use cable's rating or connected load
		current := e.CableSpec.CurrentRatingAmps
		if current == 0 {
			tgt := portToNode[e.TargetPortID]
			if tgt != nil && tgt.Type == diagram.NodeDCLoad {
				var spec diagram.DCLoadSpec
				if diagram.UnmarshalSpec(tgt.Spec, &spec) == nil && spec.Voltage > 0 {
					current = spec.Watts / spec.Voltage
				}
			}
		}

		if current > 0 && current > awg.Ampacity {
			ws = append(ws, Warning{
				Severity: SeverityError,
				EdgeID:   e.ID,
				Code:     "cable_undersized",
				Message:  fmt.Sprintf("Cable AWG %d rated for %.0fA but carrying %.0fA", e.CableSpec.GaugeAWG, awg.Ampacity, current),
			})
		}
	}
	return ws
}

func checkVoltageDrop(d *diagram.Diagram, portToNode map[string]*diagram.Node, _ map[string]*diagram.Node) []Warning {
	var ws []Warning
	for _, e := range d.Edges {
		if e.CableSpec == nil || e.CableSpec.LengthMeters == 0 {
			continue
		}
		awg, ok := LookupAWG(e.CableSpec.GaugeAWG)
		if !ok {
			continue
		}

		srcNode := portToNode[e.SourcePortID]
		tgtNode := portToNode[e.TargetPortID]

		// Find voltage from either end
		voltage := 0.0
		if srcNode != nil {
			voltage = diagram.GetNodeVoltage(srcNode)
		}
		if voltage == 0 && tgtNode != nil {
			voltage = diagram.GetNodeVoltage(tgtNode)
		}
		if voltage == 0 {
			continue
		}

		current := e.CableSpec.CurrentRatingAmps
		if current == 0 {
			continue
		}

		// V_drop = I * R * 2 (round trip)
		resistance := awg.ResistancePerMeter * e.CableSpec.LengthMeters
		vDrop := current * resistance * 2
		pct := (vDrop / voltage) * 100

		if pct > 3.0 {
			ws = append(ws, Warning{
				Severity: SeverityWarning,
				EdgeID:   e.ID,
				Code:     "voltage_drop_high",
				Message:  fmt.Sprintf("Voltage drop %.1f%% (%.2fV) on AWG %d cable, %.1fm at %.0fA — exceeds 3%% threshold", pct, vDrop, e.CableSpec.GaugeAWG, e.CableSpec.LengthMeters, current),
			})
		}
	}
	return ws
}

func checkFuseOversized(d *diagram.Diagram) []Warning {
	var ws []Warning
	for _, e := range d.Edges {
		if e.CableSpec == nil {
			continue
		}
		awg, ok := LookupAWG(e.CableSpec.GaugeAWG)
		if !ok {
			continue
		}

		// Find fuses attached to this edge's nodes (check connected fuse nodes)
		for i := range d.Nodes {
			n := &d.Nodes[i]
			if n.Type != diagram.NodeFuse {
				continue
			}
			fuseRating := diagram.GetNodeMaxAmps(n)
			if fuseRating == 0 {
				continue
			}

			// Check if this fuse is connected to the same edge
			connected := false
			for _, p := range n.Ports {
				if p.ID == e.SourcePortID || p.ID == e.TargetPortID {
					connected = true
					break
				}
			}
			if !connected {
				continue
			}

			if fuseRating > awg.Ampacity {
				ws = append(ws, Warning{
					Severity: SeverityWarning,
					NodeID:   n.ID,
					EdgeID:   e.ID,
					Code:     "fuse_oversized",
					Message:  fmt.Sprintf("Fuse %q rated %.0fA exceeds cable AWG %d ampacity (%.0fA)", n.Label, fuseRating, e.CableSpec.GaugeAWG, awg.Ampacity),
				})
			}
		}
	}
	return ws
}
