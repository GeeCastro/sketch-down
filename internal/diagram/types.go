package diagram

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const SchemaVersion = 1

// NodeType enumerates supported device types.
type NodeType string

const (
	NodeBattery             NodeType = "battery"
	NodeAlternator          NodeType = "alternator"
	NodeAlternatorRegulator NodeType = "alternatorRegulator"
	NodeMPPT                NodeType = "mppt"
	NodeSolarPanel          NodeType = "solarPanel"
	NodeFuse                NodeType = "fuse"
	NodeFuseHolder          NodeType = "fuseHolder"
	NodeBusBar              NodeType = "busBar"
	NodeDCLoad              NodeType = "dcLoad"
	NodeSwitch              NodeType = "switch"
)

type PortDirection string

const (
	PortIn            PortDirection = "in"
	PortOut           PortDirection = "out"
	PortBidirectional PortDirection = "bidirectional"
)

type PortType string

const (
	PortPositive PortType = "positive"
	PortNegative PortType = "negative"
	PortSignal   PortType = "signal"
)

type CableMaterial string

const (
	MaterialCopper       CableMaterial = "copper"
	MaterialTinnedCopper CableMaterial = "tinned_copper"
)

type FuseType string

const (
	FuseBlade  FuseType = "blade"
	FuseANL    FuseType = "ANL"
	FuseMRBF   FuseType = "MRBF"
	FuseInline FuseType = "inline"
)

type BatteryChemistry string

const (
	ChemistryAGM      BatteryChemistry = "AGM"
	ChemistryLiFePO4  BatteryChemistry = "LiFePO4"
	ChemistryLeadAcid BatteryChemistry = "lead_acid"
)

// Project is the top-level container persisted to the database.
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Diagram   *Diagram  `json:"diagram"`
}

type ViewportMeta struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

type Diagram struct {
	ID            string       `json:"id"`
	SchemaVersion int          `json:"schemaVersion"`
	Nodes         []Node       `json:"nodes"`
	Edges         []Edge       `json:"edges"`
	Metadata      ViewportMeta `json:"metadata"`
}

type Port struct {
	ID        string        `json:"id"`
	NodeID    string        `json:"nodeId"`
	Name      string        `json:"name"`
	Direction PortDirection `json:"direction"`
	PortType  PortType      `json:"portType"`
}

type Node struct {
	ID    string          `json:"id"`
	Type  NodeType        `json:"type"`
	Label string          `json:"label"`
	X     float64         `json:"x"`
	Y     float64         `json:"y"`
	Ports []Port          `json:"ports"`
	Spec  json.RawMessage `json:"spec"`
}

type CableSpec struct {
	LengthMeters     float64       `json:"lengthMeters"`
	GaugeAWG         int           `json:"gaugeAWG"`
	Material         CableMaterial `json:"material"`
	CurrentRatingAmps float64      `json:"currentRatingAmps"`
	Notes            string        `json:"notes,omitempty"`
}

type Edge struct {
	ID           string     `json:"id"`
	SourcePortID string     `json:"sourcePortId"`
	TargetPortID string     `json:"targetPortId"`
	CableSpec    *CableSpec `json:"cableSpec,omitempty"`
}

// DeviceSpec variants per NodeType.

type BatterySpec struct {
	Voltage          float64          `json:"voltage"`
	CapacityAh       float64          `json:"capacityAh"`
	Chemistry        BatteryChemistry `json:"chemistry"`
	MaxChargeAmps    float64          `json:"maxChargeAmps"`
	MaxDischargeAmps float64          `json:"maxDischargeAmps"`
}

type AlternatorSpec struct {
	Voltage       float64 `json:"voltage"`
	MaxOutputAmps float64 `json:"maxOutputAmps"`
}

type AlternatorRegulatorSpec struct {
	CompatibleChemistries []BatteryChemistry `json:"compatibleChemistries"`
}

type MPPTSpec struct {
	MaxInputVoltage float64 `json:"maxInputVoltage"`
	MaxOutputAmps   float64 `json:"maxOutputAmps"`
	Voltage         float64 `json:"voltage"`
}

type SolarPanelSpec struct {
	MaxWatts float64 `json:"maxWatts"`
	Voc      float64 `json:"voc"`
	Isc      float64 `json:"isc"`
	Vmp      float64 `json:"vmp"`
	Imp      float64 `json:"imp"`
}

type FuseSpec struct {
	RatingAmps float64  `json:"ratingAmps"`
	Type       FuseType `json:"type"`
}

type FuseHolderSpec struct {
	MaxAmps       float64    `json:"maxAmps"`
	AcceptedTypes []FuseType `json:"acceptedTypes"`
}

type BusBarSpec struct {
	MaxAmps  float64 `json:"maxAmps"`
	Circuits int     `json:"circuits"`
}

type DCLoadSpec struct {
	Watts       float64 `json:"watts"`
	Voltage     float64 `json:"voltage"`
	Description string  `json:"description"`
}

type SwitchSpec struct {
	MaxAmps float64 `json:"maxAmps"`
	Poles   int     `json:"poles"`
}

// NewID generates a new UUID string.
func NewID() string {
	return uuid.New().String()
}

// NewDiagram creates an empty diagram with defaults.
func NewDiagram() *Diagram {
	return &Diagram{
		ID:            NewID(),
		SchemaVersion: SchemaVersion,
		Nodes:         []Node{},
		Edges:         []Edge{},
		Metadata:      ViewportMeta{Zoom: 1.0},
	}
}

// NewProject creates a new project with an empty diagram.
func NewProject(name string) *Project {
	now := time.Now().UTC()
	return &Project{
		ID:        NewID(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Diagram:   NewDiagram(),
	}
}

// MarshalSpec marshals a typed spec into json.RawMessage for Node.Spec.
func MarshalSpec(spec any) (json.RawMessage, error) {
	b, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("marshal spec: %w", err)
	}
	return json.RawMessage(b), nil
}

// UnmarshalSpec unmarshals Node.Spec into a typed spec struct.
func UnmarshalSpec(raw json.RawMessage, dest any) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("unmarshal spec: %w", err)
	}
	return nil
}

// GetNodeVoltage extracts voltage from a node's spec, returning 0 if not applicable.
func GetNodeVoltage(n *Node) float64 {
	if len(n.Spec) == 0 {
		return 0
	}
	switch n.Type {
	case NodeBattery:
		var s BatterySpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.Voltage
		}
	case NodeAlternator:
		var s AlternatorSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.Voltage
		}
	case NodeMPPT:
		var s MPPTSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.Voltage
		}
	case NodeDCLoad:
		var s DCLoadSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.Voltage
		}
	}
	return 0
}

// GetNodeMaxAmps returns max current rating for current-carrying nodes.
func GetNodeMaxAmps(n *Node) float64 {
	if len(n.Spec) == 0 {
		return 0
	}
	switch n.Type {
	case NodeFuse:
		var s FuseSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.RatingAmps
		}
	case NodeSwitch:
		var s SwitchSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.MaxAmps
		}
	case NodeBusBar:
		var s BusBarSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.MaxAmps
		}
	case NodeFuseHolder:
		var s FuseHolderSpec
		if UnmarshalSpec(n.Spec, &s) == nil {
			return s.MaxAmps
		}
	}
	return 0
}
