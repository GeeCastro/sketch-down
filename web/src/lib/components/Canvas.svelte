<script lang="ts">
  import { editorState } from '../state.svelte';
  import type { Node, Port, Edge } from '../domain/types';
  import { routeEdges } from '../routing';

  function getPortType(portId: string) {
    return editorState.getNodeByPortId(portId)?.ports.find(p => p.id === portId)?.portType;
  }

  let svgElement: SVGSVGElement;

  function generateId() {
    return crypto.randomUUID();
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    const data = e.dataTransfer?.getData('application/json');
    if (!data) return;

    const nt = JSON.parse(data);
    const rect = svgElement.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const nodeId = generateId();
    const ports: Port[] = nt.ports.map((p: any) => ({
      id: generateId(),
      nodeId,
      name: p.name,
      direction: p.dir,
      portType: p.type
    }));

    const node: Node = {
      id: nodeId,
      type: nt.type,
      label: nt.label,
      imageUrl: nt.imageUrl,
      x,
      y,
      ports,
      spec: {}
    };

    editorState.addNode(node);
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'copy';
    }
  }

  // Node Dragging
  let isDraggingNode = false;
  let dragOffset = { x: 0, y: 0 };

  function startNodeDrag(e: MouseEvent, node: Node) {
    e.stopPropagation();
    editorState.selectNode(node.id);
    isDraggingNode = true;
    editorState.draggingNodeId = node.id;
    dragOffset = { x: e.clientX - node.x, y: e.clientY - node.y };
  }

  // Port Connection
  let isConnecting = false;

  function startConnection(e: MouseEvent, port: Port) {
    e.stopPropagation();
    isConnecting = true;
    editorState.connectingFromPortId = port.id;
    updateConnectionPos(e);
  }

  function updateConnectionPos(e: MouseEvent) {
    if (!isConnecting || !svgElement) return;
    const rect = svgElement.getBoundingClientRect();
    editorState.connectingToPos = {
      x: e.clientX - rect.left,
      y: e.clientY - rect.top
    };
  }

  function finishConnection(e: MouseEvent, port: Port) {
    e.stopPropagation();
    if (isConnecting && editorState.connectingFromPortId && editorState.connectingFromPortId !== port.id) {
      const edge: Edge = {
        id: generateId(),
        sourcePortId: editorState.connectingFromPortId,
        targetPortId: port.id,
        cableSpec: {
          lengthMeters: 1,
          gaugeAWG: 10,
          material: 'copper',
          currentRatingAmps: 30
        }
      };
      editorState.addEdge(edge);
    }
    cancelConnection();
  }

  function cancelConnection() {
    isConnecting = false;
    editorState.connectingFromPortId = null;
    editorState.connectingToPos = null;
  }

  function handleMouseMove(e: MouseEvent) {
    if (isDraggingNode && editorState.draggingNodeId) {
      const x = e.clientX - dragOffset.x;
      const y = e.clientY - dragOffset.y;
      editorState.updateNode(editorState.draggingNodeId, { x, y });
    } else if (isConnecting) {
      updateConnectionPos(e);
    }
  }

  function handleMouseUp() {
    if (isDraggingNode) {
      isDraggingNode = false;
      editorState.draggingNodeId = null;
      editorState.scheduleSave();
    }
    if (isConnecting) {
      cancelConnection();
    }
  }

  function handleCanvasClick() {
    editorState.selectNode(null);
    editorState.selectEdge(null);
  }

  function getPortPos(node: Node, portIndex: number, totalPorts: number) {
    const width = 100;
    const height = 60;
    // Simple layout: space ports evenly along the bottom edge for 'out', top edge for 'in'
    // For now, just space them evenly on left/right based on direction
    const port = node.ports[portIndex];
    let px = node.x;
    let py = node.y + (height / (totalPorts + 1)) * (portIndex + 1);

    if (port.direction === 'in') {
      px = node.x;
    } else if (port.direction === 'out') {
      px = node.x + width;
    } else {
      px = node.x + width / 2;
      py = node.y + height; // bottom for bidir
    }

    return { x: px, y: py };
  }

  function getPortAbsolutePos(portId: string) {
    const node = editorState.getNodeByPortId(portId);
    if (!node) return { x: 0, y: 0 };
    const idx = node.ports.findIndex(p => p.id === portId);
    return getPortPos(node, idx, node.ports.length);
  }
</script>

<svelte:window onmousemove={handleMouseMove} onmouseup={handleMouseUp} />

<svg
  bind:this={svgElement}
  class="w-full h-full"
  ondrop={handleDrop}
  ondragover={handleDragOver}
  onclick={handleCanvasClick}
  role="presentation"
>
  <defs>
    <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
      <circle cx="2" cy="2" r="1" fill="#e5e7eb" />
    </pattern>
  </defs>
  <rect width="100%" height="100%" fill="url(#grid)" />

  {#if editorState.diagram}
    {@const routedEdges = routeEdges(editorState.diagram.edges, getPortAbsolutePos, getPortType, editorState.selectedEdgeId)}
    <!-- Edges -->
    {#each routedEdges as edge (edge.id)}
      <path
        role="presentation"
        d={edge.path}
        fill="none"
        stroke={edge.color}
        stroke-width="4"
        class="cursor-pointer hover:stroke-blue-400"
        onclick={(e) => { e.stopPropagation(); editorState.selectEdge(edge.id); }}
      />
    {/each}

    <!-- Connecting line -->
    {#if isConnecting && editorState.connectingFromPortId && editorState.connectingToPos}
      {@const p1 = getPortAbsolutePos(editorState.connectingFromPortId)}
      {@const p2 = editorState.connectingToPos}
      {@const midX = (p1.x + p2.x) / 2}
      <path
        role="presentation"
        d="M {p1.x} {p1.y} L {midX} {p1.y} L {midX} {p2.y} L {p2.x} {p2.y}"
        fill="none"
        stroke="#3b82f6"
        stroke-width="2"
        stroke-dasharray="5,5"
      />
    {/if}

    <!-- Nodes -->
    {#each editorState.diagram.nodes as node (node.id)}
      <g transform="translate({node.x}, {node.y})" onclick={(e) => e.stopPropagation()}>
        <!-- Node Box -->
        <rect
          role="presentation"
          width="100"
          height="60"
          rx="6"
          fill={editorState.selectedNodeId === node.id ? '#eff6ff' : '#ffffff'}
          stroke={editorState.selectedNodeId === node.id ? '#3b82f6' : '#d1d5db'}
          stroke-width="2"
          class="cursor-move"
          onmousedown={(e) => startNodeDrag(e, node)}
        />
        {#if node.imageUrl}
          <image href={node.imageUrl} x="10" y="10" width="40" height="40" preserveAspectRatio="xMidYMid meet" class="pointer-events-none" />
          <text x="75" y="30" text-anchor="middle" dominant-baseline="middle" class="text-[10px] font-medium pointer-events-none fill-gray-700">
            {node.label.length > 10 ? node.label.substring(0, 10) + '...' : node.label}
          </text>
        {:else}
          <text x="50" y="30" text-anchor="middle" dominant-baseline="middle" class="text-xs font-medium pointer-events-none fill-gray-700">
            {node.label}
          </text>
          <text x="50" y="45" text-anchor="middle" dominant-baseline="middle" class="text-[10px] pointer-events-none fill-gray-500">
            {node.type}
          </text>
        {/if}

        <!-- Ports -->
        {#each node.ports as port, i}
          {@const pos = getPortPos(node, i, node.ports.length)}
          <circle
            cx={pos.x - node.x}
            cy={pos.y - node.y}
            r="6"
            fill={port.portType === 'positive' ? '#ef4444' : port.portType === 'negative' ? '#1f2937' : '#10b981'}
            stroke="#ffffff"
            stroke-width="2"
            class="cursor-crosshair hover:r-8 transition-all"
            onmousedown={(e) => startConnection(e, port)}
            onmouseup={(e) => finishConnection(e, port)}
            role="presentation"
          >
            <title>{port.name}</title>
          </circle>
        {/each}
      </g>
    {/each}
  {/if}
</svg>
