package diagram

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewProjectDefaults(t *testing.T) {
	p := NewProject("Test Boat")
	if p.Name != "Test Boat" {
		t.Fatalf("expected name 'Test Boat', got %q", p.Name)
	}
	if p.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if p.Diagram == nil {
		t.Fatal("expected non-nil diagram")
	}
	if p.Diagram.SchemaVersion != SchemaVersion {
		t.Fatalf("expected schema version %d, got %d", SchemaVersion, p.Diagram.SchemaVersion)
	}
	if p.Diagram.Metadata.Zoom != 1.0 {
		t.Fatalf("expected default zoom 1.0, got %f", p.Diagram.Metadata.Zoom)
	}
}

func TestProjectMarshalRoundTrip(t *testing.T) {
	p := NewProject("Round Trip")
	p.CreatedAt = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	p.UpdatedAt = p.CreatedAt

	battSpec, err := MarshalSpec(BatterySpec{
		Voltage:          12.8,
		CapacityAh:       200,
		Chemistry:        ChemistryLiFePO4,
		MaxChargeAmps:    50,
		MaxDischargeAmps: 100,
	})
	if err != nil {
		t.Fatalf("marshal battery spec: %v", err)
	}

	node := Node{
		ID:    NewID(),
		Type:  NodeBattery,
		Label: "House Bank",
		X:     100, Y: 200,
		Ports: []Port{
			{ID: NewID(), Name: "pos", Direction: PortOut, PortType: PortPositive},
			{ID: NewID(), Name: "neg", Direction: PortOut, PortType: PortNegative},
		},
		Spec: battSpec,
	}
	node.Ports[0].NodeID = node.ID
	node.Ports[1].NodeID = node.ID
	p.Diagram.Nodes = append(p.Diagram.Nodes, node)

	loadSpec, _ := MarshalSpec(DCLoadSpec{Watts: 60, Voltage: 12, Description: "Fridge"})
	loadNode := Node{
		ID:    NewID(),
		Type:  NodeDCLoad,
		Label: "Fridge",
		X:     300, Y: 200,
		Ports: []Port{
			{ID: NewID(), Name: "in+", Direction: PortIn, PortType: PortPositive},
			{ID: NewID(), Name: "in-", Direction: PortIn, PortType: PortNegative},
		},
		Spec: loadSpec,
	}
	loadNode.Ports[0].NodeID = loadNode.ID
	loadNode.Ports[1].NodeID = loadNode.ID
	p.Diagram.Nodes = append(p.Diagram.Nodes, loadNode)

	edge := Edge{
		ID:           NewID(),
		SourcePortID: node.Ports[0].ID,
		TargetPortID: loadNode.Ports[0].ID,
		CableSpec: &CableSpec{
			LengthMeters:      2.5,
			GaugeAWG:          10,
			Material:          MaterialTinnedCopper,
			CurrentRatingAmps: 30,
		},
	}
	p.Diagram.Edges = append(p.Diagram.Edges, edge)

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatalf("marshal project: %v", err)
	}

	var decoded Project
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal project: %v", err)
	}

	if decoded.Name != p.Name {
		t.Errorf("name mismatch: %q vs %q", decoded.Name, p.Name)
	}
	if len(decoded.Diagram.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(decoded.Diagram.Nodes))
	}
	if len(decoded.Diagram.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(decoded.Diagram.Edges))
	}

	var bs BatterySpec
	if err := UnmarshalSpec(decoded.Diagram.Nodes[0].Spec, &bs); err != nil {
		t.Fatalf("unmarshal battery spec: %v", err)
	}
	if bs.Voltage != 12.8 {
		t.Errorf("expected voltage 12.8, got %f", bs.Voltage)
	}
	if bs.Chemistry != ChemistryLiFePO4 {
		t.Errorf("expected LiFePO4, got %s", bs.Chemistry)
	}

	if decoded.Diagram.Edges[0].CableSpec == nil {
		t.Fatal("expected cable spec")
	}
	if decoded.Diagram.Edges[0].CableSpec.GaugeAWG != 10 {
		t.Errorf("expected gauge 10, got %d", decoded.Diagram.Edges[0].CableSpec.GaugeAWG)
	}
}

func TestMarshalSpecAllTypes(t *testing.T) {
	specs := []struct {
		name string
		spec any
	}{
		{"battery", BatterySpec{Voltage: 12.8, CapacityAh: 100, Chemistry: ChemistryAGM}},
		{"alternator", AlternatorSpec{Voltage: 14.2, MaxOutputAmps: 80}},
		{"alternatorRegulator", AlternatorRegulatorSpec{CompatibleChemistries: []BatteryChemistry{ChemistryLiFePO4}}},
		{"mppt", MPPTSpec{MaxInputVoltage: 150, MaxOutputAmps: 40, Voltage: 12}},
		{"solarPanel", SolarPanelSpec{MaxWatts: 200, Voc: 22.5, Isc: 11.1, Vmp: 18.9, Imp: 10.58}},
		{"fuse", FuseSpec{RatingAmps: 30, Type: FuseBlade}},
		{"fuseHolder", FuseHolderSpec{MaxAmps: 60, AcceptedTypes: []FuseType{FuseBlade, FuseANL}}},
		{"busBar", BusBarSpec{MaxAmps: 150, Circuits: 6}},
		{"dcLoad", DCLoadSpec{Watts: 60, Voltage: 12, Description: "Lights"}},
		{"switch", SwitchSpec{MaxAmps: 50, Poles: 2}},
	}

	for _, tc := range specs {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := MarshalSpec(tc.spec)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if len(raw) == 0 {
				t.Fatal("empty raw message")
			}

			roundTripped, err := json.Marshal(raw)
			if err != nil {
				t.Fatalf("re-marshal: %v", err)
			}
			if len(roundTripped) == 0 {
				t.Fatal("empty round-tripped data")
			}
		})
	}
}

func TestGetNodeVoltage(t *testing.T) {
	spec, _ := MarshalSpec(BatterySpec{Voltage: 12.8})
	n := &Node{Type: NodeBattery, Spec: spec}
	if v := GetNodeVoltage(n); v != 12.8 {
		t.Errorf("expected 12.8, got %f", v)
	}

	n2 := &Node{Type: NodeFuse}
	if v := GetNodeVoltage(n2); v != 0 {
		t.Errorf("expected 0 for fuse, got %f", v)
	}
}

func TestGetNodeMaxAmps(t *testing.T) {
	spec, _ := MarshalSpec(FuseSpec{RatingAmps: 30})
	n := &Node{Type: NodeFuse, Spec: spec}
	if a := GetNodeMaxAmps(n); a != 30 {
		t.Errorf("expected 30, got %f", a)
	}
}

func TestUnmarshalSpecNilSafe(t *testing.T) {
	var bs BatterySpec
	if err := UnmarshalSpec(nil, &bs); err != nil {
		t.Fatalf("unexpected error on nil: %v", err)
	}
}
