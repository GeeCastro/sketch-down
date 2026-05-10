import type { Edge, PortType } from './domain/types';

export interface Point {
  x: number;
  y: number;
}

export interface Segment {
  start: Point;
  end: Point;
  isHorizontal: boolean;
  edgeId: string;
  portType: PortType;
}

export interface RoutedEdge {
  id: string;
  path: string;
  color: string;
}

const BRIDGE_RADIUS = 6;

export function routeEdges(
  edges: Edge[],
  getPortPos: (portId: string) => Point,
  getPortType: (portId: string) => PortType | undefined,
  selectedEdgeId: string | null
): RoutedEdge[] {
  const segments: Segment[] = [];
  const edgeSegmentsMap = new Map<string, Segment[]>();
  
  // 1. Generate segments for all edges
  for (const edge of edges) {
    const p1 = getPortPos(edge.sourcePortId);
    const p2 = getPortPos(edge.targetPortId);
    const portType = getPortType(edge.sourcePortId) || 'signal';
    
    const midX = (p1.x + p2.x) / 2;
    
    const edgeSegments: Segment[] = [
      { start: p1, end: { x: midX, y: p1.y }, isHorizontal: true, edgeId: edge.id, portType },
      { start: { x: midX, y: p1.y }, end: { x: midX, y: p2.y }, isHorizontal: false, edgeId: edge.id, portType },
      { start: { x: midX, y: p2.y }, end: p2, isHorizontal: true, edgeId: edge.id, portType }
    ];
    
    segments.push(...edgeSegments);
    edgeSegmentsMap.set(edge.id, edgeSegments);
  }
  
  // 2. Find intersections
  interface Intersection {
    x: number;
    y: number;
  }
  
  const jumpsBySegment = new Map<Segment, Intersection[]>();
  
  for (let i = 0; i < segments.length; i++) {
    for (let j = i + 1; j < segments.length; j++) {
      const segA = segments[i];
      const segB = segments[j];
      
      // Do not intersect segments of the same edge
      if (segA.edgeId === segB.edgeId) continue;
      
      // We only care about orthogonal intersections (horizontal vs vertical)
      if (segA.isHorizontal === segB.isHorizontal) continue;
      
      const horiz = segA.isHorizontal ? segA : segB;
      const vert = segA.isHorizontal ? segB : segA;
      
      const hMinX = Math.min(horiz.start.x, horiz.end.x);
      const hMaxX = Math.max(horiz.start.x, horiz.end.x);
      const vMinY = Math.min(vert.start.y, vert.end.y);
      const vMaxY = Math.max(vert.start.y, vert.end.y);
      
      const ix = vert.start.x;
      const iy = horiz.start.y;
      
      if (ix > hMinX && ix < hMaxX && iy > vMinY && iy < vMaxY) {
        // Intersects! Determine who jumps.
        let aJumps = false;
        
        if (segA.portType === 'positive' && segB.portType !== 'positive') {
          aJumps = true;
        } else if (segB.portType === 'positive' && segA.portType !== 'positive') {
          aJumps = false;
        } else if (segA.portType === 'negative' && segB.portType !== 'negative') {
          aJumps = false;
        } else if (segB.portType === 'negative' && segA.portType !== 'negative') {
          aJumps = true;
        } else {
          // Default: horizontal jumps
          aJumps = segA.isHorizontal;
        }
        
        const jumper = aJumps ? segA : segB;
        if (!jumpsBySegment.has(jumper)) {
          jumpsBySegment.set(jumper, []);
        }
        jumpsBySegment.get(jumper)!.push({ x: ix, y: iy });
      }
    }
  }
  
  // 3. Build SVG paths
  const routed: RoutedEdge[] = [];
  
  for (const edge of edges) {
    const edgeSegments = edgeSegmentsMap.get(edge.id)!;
    const portType = getPortType(edge.sourcePortId) || 'signal';
    
    let path = '';
    
    for (let sIdx = 0; sIdx < edgeSegments.length; sIdx++) {
      const seg = edgeSegments[sIdx];
      const jumps = jumpsBySegment.get(seg) || [];
      
      if (sIdx === 0) {
        path += `M ${seg.start.x} ${seg.start.y} `;
      }
      
      if (jumps.length === 0) {
        path += `L ${seg.end.x} ${seg.end.y} `;
      } else {
        // Sort jumps by distance from start
        if (seg.isHorizontal) {
          const dir = seg.end.x > seg.start.x ? 1 : -1;
          jumps.sort((a, b) => (a.x - b.x) * dir);
          
          let curX = seg.start.x;
          for (const jump of jumps) {
            const jumpStartX = jump.x - BRIDGE_RADIUS * dir;
            const jumpEndX = jump.x + BRIDGE_RADIUS * dir;
            // Draw line to start of jump
            path += `L ${jumpStartX} ${seg.start.y} `;
            // Draw arc over. A rx ry x-axis-rotation large-arc-flag sweep-flag x y
            // sweep-flag is 1 if drawing in positive angle direction.
            // We want the bridge to go "up" (smaller y) for horizontal lines going right, 
            // wait, just consistent visual arc.
            // To arc "up" when going right (dir=1), sweep=1 since we go from jumpStartX to jumpEndX with negative dy
            // Actually A rx ry 0 0 1 is standard.
            // Let's just use a relative arc or absolute arc.
            const sweep = dir === 1 ? 1 : 0;
            path += `A ${BRIDGE_RADIUS} ${BRIDGE_RADIUS} 0 0 ${sweep} ${jumpEndX} ${seg.start.y} `;
            curX = jumpEndX;
          }
          path += `L ${seg.end.x} ${seg.end.y} `;
        } else {
          // Vertical segment
          const dir = seg.end.y > seg.start.y ? 1 : -1;
          jumps.sort((a, b) => (a.y - b.y) * dir);
          
          let curY = seg.start.y;
          for (const jump of jumps) {
            const jumpStartY = jump.y - BRIDGE_RADIUS * dir;
            const jumpEndY = jump.y + BRIDGE_RADIUS * dir;
            path += `L ${seg.start.x} ${jumpStartY} `;
            // sweep flag: going down (dir=1), arc to the right (larger x) -> sweep=0
            const sweep = dir === 1 ? 0 : 1;
            path += `A ${BRIDGE_RADIUS} ${BRIDGE_RADIUS} 0 0 ${sweep} ${seg.start.x} ${jumpEndY} `;
            curY = jumpEndY;
          }
          path += `L ${seg.end.x} ${seg.end.y} `;
        }
      }
    }
    
    // Determine color
    let color = '#9ca3af'; // default gray
    if (selectedEdgeId === edge.id) {
      color = '#3b82f6'; // blue
    } else {
      if (portType === 'positive') color = '#ef4444';
      else if (portType === 'negative') color = '#1f2937';
      else if (portType === 'signal') color = '#10b981';
    }
    
    routed.push({
      id: edge.id,
      path: path.trim(),
      color
    });
  }
  
  return routed;
}
