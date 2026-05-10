<script lang="ts">
  import Palette from './Palette.svelte';
  import Canvas from './Canvas.svelte';
  import Inspector from './Inspector.svelte';
  import WarningsPanel from './WarningsPanel.svelte';
  import { editorState } from '../state.svelte';
  import { getProject, exportProject, importProject } from '../api';
  import { onMount } from 'svelte';

  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      const project = await getProject();
      editorState.init(project);
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  async function handleExport() {
    try {
      await exportProject();
    } catch (e: any) {
      alert('Export failed: ' + e.message);
    }
  }

  async function handleImport(e: Event) {
    const input = e.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    try {
      const project = await importProject(input.files[0]);
      editorState.init(project);
    } catch (err: any) {
      alert('Import failed: ' + err.message);
    }
    input.value = '';
  }
</script>

<div class="flex flex-col h-screen w-screen overflow-hidden bg-gray-50 text-gray-900">
  <header class="flex items-center justify-between px-4 py-2 bg-white border-b border-gray-200 shadow-sm">
    <h1 class="text-xl font-semibold">Boat Schematics</h1>
    <div class="flex gap-2">
      <button class="px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded border border-gray-300" onclick={handleExport}>Export JSON</button>
      <label class="px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded border border-gray-300 cursor-pointer">
        Import JSON
        <input type="file" accept=".json" class="hidden" onchange={handleImport} />
      </label>
    </div>
  </header>

  {#if loading}
    <div class="flex-1 flex items-center justify-center">Loading...</div>
  {:else if error}
    <div class="flex-1 flex items-center justify-center text-red-500">{error}</div>
  {:else}
    <div class="flex flex-1 overflow-hidden">
      <!-- Left Sidebar -->
      <aside class="w-64 flex flex-col border-r border-gray-200 bg-white">
        <Palette />
      </aside>

      <!-- Main Canvas -->
      <main class="flex-1 relative overflow-hidden bg-gray-50">
        <Canvas />
      </main>

      <!-- Right Sidebar -->
      <aside class="w-80 flex flex-col border-l border-gray-200 bg-white">
        <div class="flex-1 overflow-y-auto border-b border-gray-200 p-4">
          <Inspector />
        </div>
        <div class="h-1/3 overflow-y-auto p-4 bg-gray-50">
          <WarningsPanel />
        </div>
      </aside>
    </div>
  {/if}
</div>
