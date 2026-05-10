<script lang="ts">
  import { editorState } from '../state.svelte';

  let warnings = $derived(editorState.warnings);

  function handleWarningClick(w: any) {
    if (w.nodeId) {
      editorState.selectNode(w.nodeId);
    } else if (w.edgeId) {
      editorState.selectEdge(w.edgeId);
    }
  }
</script>

<div class="h-full flex flex-col">
  <div class="flex items-center justify-between mb-3">
    <h2 class="font-semibold text-gray-700">Warnings</h2>
    <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-gray-200 text-gray-600">
      {warnings.length}
    </span>
  </div>

  <div class="flex-1 overflow-y-auto flex flex-col gap-2">
    {#if warnings.length === 0}
      <div class="text-sm text-gray-500 italic">No warnings. Diagram looks good!</div>
    {:else}
      {#each warnings as w}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
          class="p-2 text-sm rounded border-l-4 cursor-pointer transition-colors
            {w.severity === 'error' ? 'bg-red-50 border-red-500 hover:bg-red-100 text-red-800' :
             w.severity === 'warning' ? 'bg-yellow-50 border-yellow-500 hover:bg-yellow-100 text-yellow-800' :
             'bg-blue-50 border-blue-500 hover:bg-blue-100 text-blue-800'}"
          onclick={() => handleWarningClick(w)}
        >
          <div class="font-medium mb-0.5">{w.code}</div>
          <div class="opacity-90">{w.message}</div>
        </div>
      {/each}
    {/if}
  </div>
</div>
