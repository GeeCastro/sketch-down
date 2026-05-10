import type { Diagram, Project, Warning } from './domain/types';

const API_BASE = import.meta.env.DEV ? 'http://localhost:8080/api' : '/api';

export async function getProject(): Promise<Project> {
  const res = await fetch(`${API_BASE}/project`);
  if (!res.ok) throw new Error('Failed to fetch project');
  return res.json();
}

export async function saveDiagram(diagram: Diagram): Promise<void> {
  const res = await fetch(`${API_BASE}/project/diagram`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(diagram),
  });
  if (!res.ok) throw new Error('Failed to save diagram');
}

export async function validateDiagram(): Promise<Warning[]> {
  const res = await fetch(`${API_BASE}/project/validate`);
  if (!res.ok) throw new Error('Failed to validate diagram');
  const data = await res.json();
  return data.warnings;
}

export async function exportProject(): Promise<void> {
  const res = await fetch(`${API_BASE}/project/export`, { method: 'POST' });
  if (!res.ok) throw new Error('Failed to export project');
  
  const blob = await res.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  
  const contentDisposition = res.headers.get('Content-Disposition');
  let filename = 'project.json';
  if (contentDisposition) {
    const match = contentDisposition.match(/filename="?([^"]+)"?/);
    if (match) filename = match[1];
  }
  
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
}

export async function importProject(file: File): Promise<Project> {
  const res = await fetch(`${API_BASE}/project/import`, {
    method: 'POST',
    body: file, // Send file directly as body
  });
  if (!res.ok) throw new Error('Failed to import project');
  return res.json();
}
