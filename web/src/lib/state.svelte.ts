import type { Diagram, Node, Edge, Warning, Project } from './domain/types';
import { saveDiagram, validateDiagram } from './api';

export class EditorState {
  project = $state<Project | null>(null);
  diagram = $state<Diagram | null>(null);
  warnings = $state<Warning[]>([]);
  
  selectedNodeId = $state<string | null>(null);
  selectedEdgeId = $state<string | null>(null);

  // Dragging state
  draggingNodeId = $state<string | null>(null);
  
  // Connection state
  connectingFromPortId = $state<string | null>(null);
  connectingToPos = $state<{ x: number; y: number } | null>(null);

  private saveTimeout: number | null = null;

  init(project: Project) {
    this.project = project;
    this.diagram = project.diagram;
    this.validate();
  }

  get selectedNode(): Node | undefined {
    return this.diagram?.nodes.find(n => n.id === this.selectedNodeId);
  }

  get selectedEdge(): Edge | undefined {
    return this.diagram?.edges.find(e => e.id === this.selectedEdgeId);
  }

  selectNode(id: string | null) {
    this.selectedNodeId = id;
    this.selectedEdgeId = null;
  }

  selectEdge(id: string | null) {
    this.selectedEdgeId = id;
    this.selectedNodeId = null;
  }

  addNode(node: Node) {
    if (!this.diagram) return;
    this.diagram.nodes.push(node);
    this.scheduleSave();
  }

  updateNode(id: string, updates: Partial<Node>) {
    if (!this.diagram) return;
    const node = this.diagram.nodes.find(n => n.id === id);
    if (node) {
      Object.assign(node, updates);
      this.scheduleSave();
    }
  }

  addEdge(edge: Edge) {
    if (!this.diagram) return;
    this.diagram.edges.push(edge);
    this.scheduleSave();
  }

  updateEdge(id: string, updates: Partial<Edge>) {
    if (!this.diagram) return;
    const edge = this.diagram.edges.find(e => e.id === id);
    if (edge) {
      Object.assign(edge, updates);
      this.scheduleSave();
    }
  }

  deleteSelected() {
    if (!this.diagram) return;
    if (this.selectedNodeId) {
      this.diagram.nodes = this.diagram.nodes.filter(n => n.id !== this.selectedNodeId);
      // Also delete connected edges
      this.diagram.edges = this.diagram.edges.filter(e => {
        const sourceNode = this.getNodeByPortId(e.sourcePortId);
        const targetNode = this.getNodeByPortId(e.targetPortId);
        return sourceNode?.id !== this.selectedNodeId && targetNode?.id !== this.selectedNodeId;
      });
      this.selectedNodeId = null;
      this.scheduleSave();
    } else if (this.selectedEdgeId) {
      this.diagram.edges = this.diagram.edges.filter(e => e.id !== this.selectedEdgeId);
      this.selectedEdgeId = null;
      this.scheduleSave();
    }
  }

  getNodeByPortId(portId: string): Node | undefined {
    return this.diagram?.nodes.find(n => n.ports.some(p => p.id === portId));
  }

  scheduleSave() {
    if (this.saveTimeout) {
      clearTimeout(this.saveTimeout);
    }
    this.saveTimeout = window.setTimeout(async () => {
      if (this.diagram) {
        try {
          await saveDiagram(this.diagram);
          await this.validate();
        } catch (e) {
          console.error("Failed to save diagram", e);
        }
      }
    }, 500);
  }

  async validate() {
    try {
      this.warnings = await validateDiagram();
    } catch (e) {
      console.error("Failed to validate diagram", e);
    }
  }
}

export const editorState = new EditorState();
