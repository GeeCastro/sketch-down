<script lang="ts">
  import { editorState } from '../state.svelte';
  import type { Node, Edge, CableSpec } from '../domain/types';

  let selectedNode = $derived(editorState.selectedNode);
  let selectedEdge = $derived(editorState.selectedEdge);

  function handleNodeChange(e: Event, field: string) {
    if (!selectedNode) return;
    const target = e.target as HTMLInputElement;
    const val = target.type === 'number' ? parseFloat(target.value) : target.value;
    
    if (field === 'label') {
      editorState.updateNode(selectedNode.id, { label: val as string });
    } else {
      const spec = { ...selectedNode.spec, [field]: val };
      editorState.updateNode(selectedNode.id, { spec });
    }
  }

  function handleEdgeChange(e: Event, field: keyof CableSpec) {
    if (!selectedEdge || !selectedEdge.cableSpec) return;
    const target = e.target as HTMLInputElement;
    const val = target.type === 'number' ? parseFloat(target.value) : target.value;
    
    const cableSpec = { ...selectedEdge.cableSpec, [field]: val };
    editorState.updateEdge(selectedEdge.id, { cableSpec });
  }

  function handleDelete() {
    editorState.deleteSelected();
  }
</script>

<div class="h-full">
  <h2 class="font-semibold mb-4 text-gray-700">Inspector</h2>

  {#if selectedNode}
    <div class="flex flex-col gap-4">
      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1" for="node-label">Label</label>
        <input
          id="node-label"
          type="text"
          value={selectedNode.label}
          oninput={(e) => handleNodeChange(e, 'label')}
          class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none"
        />
      </div>

      <!-- Basic Spec Fields based on type -->
      {#if ['battery', 'alternator', 'mppt', 'dcLoad'].includes(selectedNode.type)}
        <div>
          <label class="block text-xs font-medium text-gray-500 mb-1" for="node-voltage">Voltage (V)</label>
          <input
            id="node-voltage"
            type="number"
            value={selectedNode.spec?.voltage || ''}
            oninput={(e) => handleNodeChange(e, 'voltage')}
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
          />
        </div>
      {/if}

      {#if selectedNode.type === 'battery'}
        <div>
          <label class="block text-xs font-medium text-gray-500 mb-1" for="node-capacity">Capacity (Ah)</label>
          <input
            id="node-capacity"
            type="number"
            value={selectedNode.spec?.capacityAh || ''}
            oninput={(e) => handleNodeChange(e, 'capacityAh')}
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
          />
        </div>
      {/if}

      {#if selectedNode.type === 'fuse'}
        <div>
          <label class="block text-xs font-medium text-gray-500 mb-1" for="node-rating">Rating (Amps)</label>
          <input
            id="node-rating"
            type="number"
            value={selectedNode.spec?.ratingAmps || ''}
            oninput={(e) => handleNodeChange(e, 'ratingAmps')}
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
          />
        </div>
      {/if}

      {#if selectedNode.type === 'dcLoad'}
        <div>
          <label class="block text-xs font-medium text-gray-500 mb-1" for="node-watts">Watts</label>
          <input
            id="node-watts"
            type="number"
            value={selectedNode.spec?.watts || ''}
            oninput={(e) => handleNodeChange(e, 'watts')}
            class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
          />
        </div>
      {/if}

      <button onclick={handleDelete} class="mt-4 px-3 py-1.5 text-sm text-red-600 bg-red-50 hover:bg-red-100 rounded border border-red-200">
        Delete Node
      </button>
    </div>

  {:else if selectedEdge}
    <div class="flex flex-col gap-4">
      <div class="text-sm font-medium text-gray-700">Cable</div>
      
      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1" for="edge-length">Length (Meters)</label>
        <input
          id="edge-length"
          type="number"
          value={selectedEdge.cableSpec?.lengthMeters || ''}
          oninput={(e) => handleEdgeChange(e, 'lengthMeters')}
          class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
        />
      </div>

      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1" for="edge-gauge">Gauge (AWG)</label>
        <input
          id="edge-gauge"
          type="number"
          value={selectedEdge.cableSpec?.gaugeAWG || ''}
          oninput={(e) => handleEdgeChange(e, 'gaugeAWG')}
          class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
        />
      </div>

      <div>
        <label class="block text-xs font-medium text-gray-500 mb-1" for="edge-current">Current Rating (Amps)</label>
        <input
          id="edge-current"
          type="number"
          value={selectedEdge.cableSpec?.currentRatingAmps || ''}
          oninput={(e) => handleEdgeChange(e, 'currentRatingAmps')}
          class="w-full px-2 py-1 text-sm border border-gray-300 rounded outline-none"
        />
      </div>

      <button onclick={handleDelete} class="mt-4 px-3 py-1.5 text-sm text-red-600 bg-red-50 hover:bg-red-100 rounded border border-red-200">
        Delete Cable
      </button>
    </div>

  {:else}
    <div class="text-sm text-gray-500 italic">Select a node or cable to edit its properties.</div>
  {/if}
</div>
