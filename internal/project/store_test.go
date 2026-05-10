package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sketch-down/sketch-down/internal/diagram"
)

func tempStore(t *testing.T) Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestCreateAndGetActive(t *testing.T) {
	s := tempStore(t)

	p, err := s.CreateProject("My Boat")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.Name != "My Boat" {
		t.Errorf("name: got %q", p.Name)
	}

	active, err := s.GetActiveProject()
	if err != nil {
		t.Fatalf("get active: %v", err)
	}
	if active == nil {
		t.Fatal("expected active project")
	}
	if active.ID != p.ID {
		t.Errorf("active ID mismatch: %s vs %s", active.ID, p.ID)
	}
}

func TestSaveDiagram(t *testing.T) {
	s := tempStore(t)
	p, _ := s.CreateProject("Test")

	battSpec, _ := diagram.MarshalSpec(diagram.BatterySpec{Voltage: 12.8, CapacityAh: 200})
	node := diagram.Node{
		ID:    diagram.NewID(),
		Type:  diagram.NodeBattery,
		Label: "Bank 1",
		Spec:  battSpec,
	}
	p.Diagram.Nodes = append(p.Diagram.Nodes, node)

	if err := s.SaveDiagram(p.ID, p.Diagram); err != nil {
		t.Fatalf("save: %v", err)
	}

	reloaded, err := s.GetActiveProject()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(reloaded.Diagram.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(reloaded.Diagram.Nodes))
	}
	if reloaded.Diagram.Nodes[0].Label != "Bank 1" {
		t.Errorf("label mismatch: %q", reloaded.Diagram.Nodes[0].Label)
	}
}

func TestSaveDiagramNotFound(t *testing.T) {
	s := tempStore(t)
	d := diagram.NewDiagram()
	err := s.SaveDiagram("nonexistent", d)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestListProjects(t *testing.T) {
	s := tempStore(t)
	s.CreateProject("A")
	s.CreateProject("B")

	list, err := s.ListProjects()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(list))
	}
}

func TestActiveProjectSwitches(t *testing.T) {
	s := tempStore(t)
	a, _ := s.CreateProject("A")
	b, _ := s.CreateProject("B")

	active, _ := s.GetActiveProject()
	if active.ID != b.ID {
		t.Errorf("expected B active, got %s", active.ID)
	}
	_ = a // a was deactivated when b was created
}

func TestExportImport(t *testing.T) {
	s := tempStore(t)
	p, _ := s.CreateProject("Export Me")
	battSpec, _ := diagram.MarshalSpec(diagram.BatterySpec{Voltage: 24})
	p.Diagram.Nodes = append(p.Diagram.Nodes, diagram.Node{
		ID: diagram.NewID(), Type: diagram.NodeBattery, Label: "24V", Spec: battSpec,
	})
	s.SaveDiagram(p.ID, p.Diagram)

	data, err := s.ExportProject(p.ID)
	if err != nil {
		t.Fatalf("export: %v", err)
	}

	// Verify it's valid JSON
	var exported diagram.Project
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("exported JSON invalid: %v", err)
	}
	if exported.Name != "Export Me" {
		t.Errorf("exported name: %q", exported.Name)
	}

	// Import into a fresh store
	s2 := tempStore(t)
	imported, err := s2.ImportProject(data)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if imported.Name != "Export Me" {
		t.Errorf("imported name: %q", imported.Name)
	}

	active, _ := s2.GetActiveProject()
	if active == nil {
		t.Fatal("expected active after import")
	}
	if len(active.Diagram.Nodes) != 1 {
		t.Fatalf("expected 1 node after import, got %d", len(active.Diagram.Nodes))
	}
}

func TestExportNotFound(t *testing.T) {
	s := tempStore(t)
	_, err := s.ExportProject("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestImportInvalidJSON(t *testing.T) {
	s := tempStore(t)
	_, err := s.ImportProject([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGetActiveProjectEmpty(t *testing.T) {
	s := tempStore(t)
	p, err := s.GetActiveProject()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Fatal("expected nil for empty db")
	}
}

func TestImportWritesToDisk(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "sub", "test.db")
	s, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	defer s.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected db file to exist")
	}
}
