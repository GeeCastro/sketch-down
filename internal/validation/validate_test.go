package validation

import (
	"testing"

	"github.com/sketch-down/sketch-down/internal/diagram"
)

func mustSpec(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := diagram.MarshalSpec(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func batteryNode(t *testing.T, label string, voltage float64) diagram.Node {
	n := diagram.Node{
		ID:    diagram.NewID(),
		Type:  diagram.NodeBattery,
		Label: label,
		Ports: []diagram.Port{
			{ID: diagram.NewID(), Name: "pos", Direction: diagram.PortOut, PortType: diagram.PortPositive},
			{ID: diagram.NewID(), Name: "neg", Direction: diagram.PortOut, PortType: diagram.PortNegative},
		},
	}
	n.Ports[0].NodeID = n.ID
	n.Ports[1].NodeID = n.ID
	if voltage > 0 {
		n.Spec = mustSpec(t, diagram.BatterySpec{Voltage: voltage, CapacityAh: 100, Chemistry: diagram.ChemistryLiFePO4})
	}
	return n
}

func loadNode(t *testing.T, label string, watts, voltage float64) diagram.Node {
	n := diagram.Node{
		ID:    diagram.NewID(),
		Type:  diagram.NodeDCLoad,
		Label: label,
		Ports: []diagram.Port{
			{ID: diagram.NewID(), Name: "in+", Direction: diagram.PortIn, PortType: diagram.PortPositive},
			{ID: diagram.NewID(), Name: "in-", Direction: diagram.PortIn, PortType: diagram.PortNegative},
		},
	}
	n.Ports[0].NodeID = n.ID
	n.Ports[1].NodeID = n.ID
	if watts > 0 || voltage > 0 {
		n.Spec = mustSpec(t, diagram.DCLoadSpec{Watts: watts, Voltage: voltage, Description: label})
	}
	return n
}

func fuseNode(t *testing.T, label string, rating float64) diagram.Node {
	n := diagram.Node{
		ID:    diagram.NewID(),
		Type:  diagram.NodeFuse,
		Label: label,
		Ports: []diagram.Port{
			{ID: diagram.NewID(), Name: "in", Direction: diagram.PortIn, PortType: diagram.PortPositive},
			{ID: diagram.NewID(), Name: "out", Direction: diagram.PortOut, PortType: diagram.PortPositive},
		},
	}
	n.Ports[0].NodeID = n.ID
	n.Ports[1].NodeID = n.ID
	if rating > 0 {
		n.Spec = mustSpec(t, diagram.FuseSpec{RatingAmps: rating, Type: diagram.FuseBlade})
	}
	return n
}

func connect(src, tgt diagram.Node, cable *diagram.CableSpec) diagram.Edge {
	return diagram.Edge{
		ID:           diagram.NewID(),
		SourcePortID: src.Ports[0].ID,
		TargetPortID: tgt.Ports[0].ID,
		CableSpec:    cable,
	}
}

func hasCode(ws []Warning, code string) bool {
	for _, w := range ws {
		if w.Code == code {
			return true
		}
	}
	return false
}

func countCode(ws []Warning, code string) int {
	n := 0
	for _, w := range ws {
		if w.Code == code {
			n++
		}
	}
	return n
}

func TestMissingVoltage(t *testing.T) {
	batt := batteryNode(t, "No Voltage", 0) // no spec = voltage 0
	batt.Spec = nil
	d := &diagram.Diagram{Nodes: []diagram.Node{batt}}
	ws := Validate(d)
	if !hasCode(ws, "missing_voltage") {
		t.Error("expected missing_voltage warning")
	}
}

func TestMissingVoltageOK(t *testing.T) {
	batt := batteryNode(t, "Has Voltage", 12.8)
	d := &diagram.Diagram{Nodes: []diagram.Node{batt}}
	ws := Validate(d)
	if hasCode(ws, "missing_voltage") {
		t.Error("unexpected missing_voltage warning")
	}
}

func TestMissingFuseRating(t *testing.T) {
	f := fuseNode(t, "Empty Fuse", 0)
	f.Spec = nil
	d := &diagram.Diagram{Nodes: []diagram.Node{f}}
	ws := Validate(d)
	if !hasCode(ws, "missing_fuse_rating") {
		t.Error("expected missing_fuse_rating warning")
	}
}

func TestMissingFuseRatingOK(t *testing.T) {
	f := fuseNode(t, "30A Fuse", 30)
	d := &diagram.Diagram{Nodes: []diagram.Node{f}}
	ws := Validate(d)
	if hasCode(ws, "missing_fuse_rating") {
		t.Error("unexpected missing_fuse_rating warning")
	}
}

func TestVoltageMismatch(t *testing.T) {
	batt := batteryNode(t, "12V Batt", 12)
	load := loadNode(t, "24V Load", 60, 24)
	edge := connect(batt, load, nil)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if !hasCode(ws, "voltage_mismatch") {
		t.Error("expected voltage_mismatch")
	}
}

func TestVoltageMismatchOK(t *testing.T) {
	batt := batteryNode(t, "12V Batt", 12)
	load := loadNode(t, "12V Load", 60, 12)
	edge := connect(batt, load, nil)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if hasCode(ws, "voltage_mismatch") {
		t.Error("unexpected voltage_mismatch")
	}
}

func TestCableNoFuse(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	load := loadNode(t, "Fridge", 60, 12)
	edge := connect(batt, load, nil)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if !hasCode(ws, "cable_no_fuse") {
		t.Error("expected cable_no_fuse")
	}
}

func TestCableWithFuse(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	fuse := fuseNode(t, "30A", 30)
	load := loadNode(t, "Fridge", 60, 12)

	e1 := connect(batt, fuse, nil)
	e2 := diagram.Edge{
		ID:           diagram.NewID(),
		SourcePortID: fuse.Ports[1].ID,
		TargetPortID: load.Ports[0].ID,
	}

	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, fuse, load},
		Edges: []diagram.Edge{e1, e2},
	}
	ws := Validate(d)
	if hasCode(ws, "cable_no_fuse") {
		t.Error("unexpected cable_no_fuse — fuse is present")
	}
}

func TestCableUndersized(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	load := loadNode(t, "Windlass", 600, 12) // 50A
	cable := &diagram.CableSpec{GaugeAWG: 16, CurrentRatingAmps: 50, LengthMeters: 3, Material: diagram.MaterialCopper}
	edge := connect(batt, load, cable) // AWG 16 = 25A, but 50A flowing
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if !hasCode(ws, "cable_undersized") {
		t.Error("expected cable_undersized")
	}
}

func TestCableUndersizedOK(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	load := loadNode(t, "Light", 12, 12) // 1A
	cable := &diagram.CableSpec{GaugeAWG: 16, CurrentRatingAmps: 1, LengthMeters: 2, Material: diagram.MaterialCopper}
	edge := connect(batt, load, cable)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if hasCode(ws, "cable_undersized") {
		t.Error("unexpected cable_undersized")
	}
}

func TestVoltageDropHigh(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	load := loadNode(t, "Bow Thruster", 600, 12)
	// AWG 16, 25A capacity, long cable, high current => high drop
	cable := &diagram.CableSpec{GaugeAWG: 16, CurrentRatingAmps: 20, LengthMeters: 15, Material: diagram.MaterialCopper}
	edge := connect(batt, load, cable)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if !hasCode(ws, "voltage_drop_high") {
		t.Error("expected voltage_drop_high")
	}
}

func TestVoltageDropOK(t *testing.T) {
	batt := batteryNode(t, "Battery", 12)
	load := loadNode(t, "LED", 3, 12)
	cable := &diagram.CableSpec{GaugeAWG: 10, CurrentRatingAmps: 0.5, LengthMeters: 2, Material: diagram.MaterialCopper}
	edge := connect(batt, load, cable)
	d := &diagram.Diagram{
		Nodes: []diagram.Node{batt, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if hasCode(ws, "voltage_drop_high") {
		t.Error("unexpected voltage_drop_high")
	}
}

func TestFuseOversized(t *testing.T) {
	fuse := fuseNode(t, "Big Fuse", 100)
	load := loadNode(t, "Light", 12, 12)
	cable := &diagram.CableSpec{GaugeAWG: 16, CurrentRatingAmps: 10, LengthMeters: 1, Material: diagram.MaterialCopper}
	// Edge from fuse out port to load, with thin cable
	edge := diagram.Edge{
		ID:           diagram.NewID(),
		SourcePortID: fuse.Ports[1].ID,
		TargetPortID: load.Ports[0].ID,
		CableSpec:    cable,
	}
	d := &diagram.Diagram{
		Nodes: []diagram.Node{fuse, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if !hasCode(ws, "fuse_oversized") {
		t.Error("expected fuse_oversized")
	}
}

func TestFuseOversizedOK(t *testing.T) {
	fuse := fuseNode(t, "20A Fuse", 20)
	load := loadNode(t, "Light", 12, 12)
	cable := &diagram.CableSpec{GaugeAWG: 10, CurrentRatingAmps: 10, LengthMeters: 1, Material: diagram.MaterialCopper}
	edge := diagram.Edge{
		ID:           diagram.NewID(),
		SourcePortID: fuse.Ports[1].ID,
		TargetPortID: load.Ports[0].ID,
		CableSpec:    cable,
	}
	d := &diagram.Diagram{
		Nodes: []diagram.Node{fuse, load},
		Edges: []diagram.Edge{edge},
	}
	ws := Validate(d)
	if hasCode(ws, "fuse_oversized") {
		t.Error("unexpected fuse_oversized — fuse 20A < AWG 10 ampacity 60A")
	}
}

func TestValidateNilDiagram(t *testing.T) {
	ws := Validate(nil)
	if ws != nil {
		t.Error("expected nil for nil diagram")
	}
}

func TestValidateEmptyDiagram(t *testing.T) {
	ws := Validate(&diagram.Diagram{})
	if len(ws) != 0 {
		t.Errorf("expected 0 warnings for empty diagram, got %d", len(ws))
	}
}
