<script lang="ts">
  import { editorState } from '../state.svelte';
  import type { NodeType, Node, Port, PortType, PortDirection } from '../domain/types';

  const nodeTypes: { type: NodeType; label: string; imageUrl?: string; ports: { name: string; dir: PortDirection; type: PortType }[] }[] = [
    {
      type: 'sterlingCombiPro',
      label: 'Sterling Combi Pro S',
      imageUrl: '/components/sterling-combi-pro.png',
      ports: [
        { name: 'DC+', dir: 'in', type: 'positive' },
        { name: 'DC-', dir: 'in', type: 'negative' },
        { name: 'AC IN', dir: 'in', type: 'signal' },
        { name: 'AC OUT', dir: 'out', type: 'signal' }
      ]
    },
    {
      type: 'victronMultiPlus',
      label: 'Victron MultiPlus',
      imageUrl: '/components/victron-multiplus.png',
      ports: [
        { name: 'DC+', dir: 'in', type: 'positive' },
        { name: 'DC-', dir: 'in', type: 'negative' },
        { name: 'AC IN', dir: 'in', type: 'signal' },
        { name: 'AC OUT 1', dir: 'out', type: 'signal' },
        { name: 'AC OUT 2', dir: 'out', type: 'signal' },
        { name: 'VE.Bus', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'renogyInverter',
      label: 'Renogy Inverter',
      imageUrl: '/components/renogy-inverter.png',
      ports: [
        { name: 'DC+', dir: 'in', type: 'positive' },
        { name: 'DC-', dir: 'in', type: 'negative' },
        { name: 'AC OUT', dir: 'out', type: 'signal' }
      ]
    },
    {
      type: 'genericAlternator',
      label: 'Generic Alternator',
      imageUrl: '/components/generic-alternator.png',
      ports: [
        { name: 'B+', dir: 'out', type: 'positive' },
        { name: 'B-', dir: 'out', type: 'negative' },
        { name: 'Field', dir: 'in', type: 'signal' },
        { name: 'Stator', dir: 'out', type: 'signal' }
      ]
    },
    {
      type: 'wakespeedWS500',
      label: 'Wakespeed WS500 Pro',
      imageUrl: '/components/wakespeed-ws500.png',
      ports: [
        { name: 'Ampseal', dir: 'bidirectional', type: 'signal' },
        { name: 'CANbus', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'mastervoltAlphaPro',
      label: 'Mastervolt Alpha Pro III',
      imageUrl: '/components/mastervolt-alpha-pro.png',
      ports: [
        { name: 'Plug', dir: 'bidirectional', type: 'signal' },
        { name: 'MasterBus', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'fogstarDriftGen2',
      label: 'Fogstar Drift Gen 2 628Ah',
      imageUrl: '/components/fogstar-drift-gen2.png',
      ports: [
        { name: 'POS', dir: 'bidirectional', type: 'positive' },
        { name: 'NEG', dir: 'bidirectional', type: 'negative' }
      ]
    },
    {
      type: 'victronSmartSolarMPPT',
      label: 'Victron SmartSolar MPPT',
      imageUrl: '/components/victron-smartsolar-mppt.png',
      ports: [
        { name: 'PV+', dir: 'in', type: 'positive' },
        { name: 'PV-', dir: 'in', type: 'negative' },
        { name: 'BATT+', dir: 'out', type: 'positive' },
        { name: 'BATT-', dir: 'out', type: 'negative' },
        { name: 'VE.Direct', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'victronCerboGXMK2',
      label: 'Victron Cerbo GX MK2',
      imageUrl: '/components/victron-cerbogx-mk2.webp',
      ports: [
        { name: 'Power IN', dir: 'in', type: 'positive' },
        { name: 'VE.Direct', dir: 'bidirectional', type: 'signal' },
        { name: 'VE.Bus', dir: 'bidirectional', type: 'signal' },
        { name: 'VE.Can', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'starlink12V',
      label: 'Starlink 12V Conversion Kit',
      imageUrl: '/components/starlink-12v.png',
      ports: [
        { name: 'DC+', dir: 'in', type: 'positive' },
        { name: 'DC-', dir: 'in', type: 'negative' },
        { name: 'RJ45', dir: 'out', type: 'signal' }
      ]
    },
    {
      type: 'victronSmartShunt',
      label: 'Victron SmartShunt',
      imageUrl: '/components/victron-smartshunt.png',
      ports: [
        { name: 'BATT-', dir: 'in', type: 'negative' },
        { name: 'SYSTEM-', dir: 'out', type: 'negative' },
        { name: 'VE.Direct', dir: 'bidirectional', type: 'signal' }
      ]
    },
    {
      type: 'fuse',
      label: 'Fuse',
      ports: [
        { name: 'in', dir: 'in', type: 'positive' },
        { name: 'out', dir: 'out', type: 'positive' }
      ]
    },
    {
      type: 'busBar',
      label: 'Bus Bar',
      ports: [
        { name: 'p1', dir: 'bidirectional', type: 'positive' },
        { name: 'p2', dir: 'bidirectional', type: 'positive' },
        { name: 'p3', dir: 'bidirectional', type: 'positive' },
        { name: 'p4', dir: 'bidirectional', type: 'positive' }
      ]
    },
    {
      type: 'dcLoad',
      label: 'DC Load',
      ports: [
        { name: 'in+', dir: 'in', type: 'positive' },
        { name: 'in-', dir: 'in', type: 'negative' }
      ]
    },
    {
      type: 'switch',
      label: 'Switch',
      ports: [
        { name: 'in', dir: 'in', type: 'positive' },
        { name: 'out', dir: 'out', type: 'positive' }
      ]
    }
  ];

  function handleDragStart(e: DragEvent, nt: typeof nodeTypes[0]) {
    if (e.dataTransfer) {
      e.dataTransfer.setData('application/json', JSON.stringify(nt));
      e.dataTransfer.effectAllowed = 'copy';
    }
  }
</script>

<div class="p-4 h-full overflow-y-auto">
  <h2 class="font-semibold mb-4 text-gray-700">Components</h2>
  <div class="flex flex-col gap-2">
    {#each nodeTypes as nt}
      <div
        role="button"
        tabindex="0"
        draggable="true"
        ondragstart={(e) => handleDragStart(e, nt)}
        class="p-3 bg-white border border-gray-300 rounded shadow-sm cursor-grab hover:border-blue-500 hover:shadow flex items-center gap-3"
      >
        {#if nt.imageUrl}
          <img src={nt.imageUrl} alt={nt.label} class="w-8 h-8 object-contain" />
        {/if}
        <div class="font-medium text-sm leading-tight">{nt.label}</div>
      </div>
    {/each}
  </div>
</div>
