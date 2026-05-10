export type NodeType =
  | 'battery'
  | 'alternator'
  | 'alternatorRegulator'
  | 'mppt'
  | 'solarPanel'
  | 'fuse'
  | 'fuseHolder'
  | 'busBar'
  | 'dcLoad'
  | 'switch'
  | 'sterlingCombiPro'
  | 'victronMultiPlus'
  | 'renogyInverter'
  | 'genericAlternator'
  | 'wakespeedWS500'
  | 'mastervoltAlphaPro'
  | 'fogstarDriftGen2'
  | 'victronSmartSolarMPPT'
  | 'victronCerboGXMK2'
  | 'starlink12V'
  | 'victronSmartShunt';

export type PortDirection = 'in' | 'out' | 'bidirectional';
export type PortType = 'positive' | 'negative' | 'signal';
export type CableMaterial = 'copper' | 'tinned_copper';
export type FuseType = 'blade' | 'ANL' | 'MRBF' | 'inline';
export type BatteryChemistry = 'AGM' | 'LiFePO4' | 'lead_acid';

export interface Project {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  diagram: Diagram;
}

export interface ViewportMeta {
  x: number;
  y: number;
  zoom: number;
}

export interface Diagram {
  id: string;
  schemaVersion: number;
  nodes: Node[];
  edges: Edge[];
  metadata: ViewportMeta;
}

export interface Port {
  id: string;
  nodeId: string;
  name: string;
  direction: PortDirection;
  portType: PortType;
}

export interface Node {
  id: string;
  type: NodeType;
  label: string;
  x: number;
  y: number;
  ports: Port[];
  spec: any; // DeviceSpec variant
  imageUrl?: string;
}

export interface CableSpec {
  lengthMeters: number;
  gaugeAWG: number;
  material: CableMaterial;
  currentRatingAmps: number;
  notes?: string;
}

export interface Edge {
  id: string;
  sourcePortId: string;
  targetPortId: string;
  cableSpec?: CableSpec;
}

// Device Specs

export interface BatterySpec {
  voltage: number;
  capacityAh: number;
  chemistry: BatteryChemistry;
  maxChargeAmps: number;
  maxDischargeAmps: number;
}

export interface AlternatorSpec {
  voltage: number;
  maxOutputAmps: number;
}

export interface AlternatorRegulatorSpec {
  compatibleChemistries: BatteryChemistry[];
}

export interface MPPTSpec {
  maxInputVoltage: number;
  maxOutputAmps: number;
  voltage: number;
}

export interface SolarPanelSpec {
  maxWatts: number;
  voc: number;
  isc: number;
  vmp: number;
  imp: number;
}

export interface FuseSpec {
  ratingAmps: number;
  type: FuseType;
}

export interface FuseHolderSpec {
  maxAmps: number;
  acceptedTypes: FuseType[];
}

export interface BusBarSpec {
  maxAmps: number;
  circuits: number;
}

export interface DCLoadSpec {
  watts: number;
  voltage: number;
  description: string;
}

export interface SwitchSpec {
  maxAmps: number;
  poles: number;
}

// Validation

export type Severity = 'info' | 'warning' | 'error';

export interface Warning {
  severity: Severity;
  nodeId?: string;
  edgeId?: string;
  code: string;
  message: string;
}
