package go_path_tree

import (
	"reflect"
	"testing"
)

func TestNewTreeHasEmptyDirectoryRoot(t *testing.T) {
	tr := NewTree()
	if tr.Root == nil {
		t.Fatal("Root is nil")
	}
	if tr.Root.Name != "" {
		t.Errorf("Root.Name = %q, want empty", tr.Root.Name)
	}
	if !tr.Root.IsDir {
		t.Error("Root.IsDir = false, want true")
	}
	if len(tr.Root.Children) != 0 {
		t.Errorf("Root has %d children, want 0", len(tr.Root.Children))
	}
}

func TestAddSingleFile(t *testing.T) {
	tr := NewTree()
	tr.Add("a.txt")

	if len(tr.Root.Children) != 1 {
		t.Fatalf("Root has %d children, want 1", len(tr.Root.Children))
	}
	child := tr.Root.Children[0]
	if child.Name != "a.txt" {
		t.Errorf("Name = %q, want %q", child.Name, "a.txt")
	}
	if child.IsDir {
		t.Error("IsDir = true, want false")
	}
	if child.Children != nil {
		t.Errorf("file has children: %v", child.Children)
	}
}

func TestAddNestedPathCreatesDirectories(t *testing.T) {
	tr := NewTree()
	tr.Add("a/b/c.txt")

	a := tr.Root.Children[0]
	if a.Name != "a" || !a.IsDir {
		t.Fatalf("first child = %+v, want dir named a", a)
	}
	b := a.Children[0]
	if b.Name != "b" || !b.IsDir {
		t.Fatalf("second level = %+v, want dir named b", b)
	}
	c := b.Children[0]
	if c.Name != "c.txt" || c.IsDir {
		t.Fatalf("leaf = %+v, want file named c.txt", c)
	}
}

func TestAddRepeatedComponentsAreIgnored(t *testing.T) {
	tr := NewTree()
	tr.Add("a//b/")

	if len(tr.Root.Children) != 1 {
		t.Fatalf("Root has %d children, want 1", len(tr.Root.Children))
	}
	a := tr.Root.Children[0]
	if a.Name != "a" || !a.IsDir {
		t.Fatalf("a = %+v, want dir", a)
	}
	if len(a.Children) != 1 {
		t.Fatalf("a has %d children, want 1", len(a.Children))
	}
	b := a.Children[0]
	if b.Name != "b" || b.IsDir {
		t.Fatalf("b = %+v, want file", b)
	}
}

func TestAddEmptyPathIsNoOp(t *testing.T) {
	tr := NewTree()
	tr.Add("")
	tr.Add("/")
	tr.Add("///")
	if len(tr.Root.Children) != 0 {
		t.Errorf("Root has %d children, want 0", len(tr.Root.Children))
	}
}

func TestAddSamePathTwiceDoesNotDuplicate(t *testing.T) {
	tr := NewTree()
	tr.Add("a/b.txt")
	tr.Add("a/b.txt")

	if len(tr.Root.Children) != 1 {
		t.Fatalf("Root has %d children, want 1", len(tr.Root.Children))
	}
	a := tr.Root.Children[0]
	if len(a.Children) != 1 {
		t.Fatalf("a has %d children, want 1", len(a.Children))
	}
}

func TestAddFileThenDescendantStopsAtFile(t *testing.T) {
	tr := NewTree()
	tr.Add("a")
	tr.Add("a/b")

	if len(tr.Root.Children) != 1 {
		t.Fatalf("Root has %d children, want 1", len(tr.Root.Children))
	}
	a := tr.Root.Children[0]
	if a.IsDir {
		t.Fatal("a became a directory, want file")
	}
	if len(a.Children) != 0 {
		t.Errorf("file a gained %d children", len(a.Children))
	}
}

func TestAddAllInsertsInOrder(t *testing.T) {
	tr := NewTree()
	tr.AddAll([]string{"b.txt", "a/c.txt", "a/d.txt"})

	paths := tr.Paths()
	want := []string{"b.txt", "a/c.txt", "a/d.txt"}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("Paths() = %v, want %v", paths, want)
	}
}

func TestSortOrdersChildrenByName(t *testing.T) {
	tr := NewTree()
	tr.AddAll([]string{"z.txt", "a.txt", "m/b.txt", "m/a.txt"})
	tr.Sort()

	paths := tr.Paths()
	want := []string{"a.txt", "m/a.txt", "m/b.txt", "z.txt"}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("after Sort, Paths() = %v, want %v", paths, want)
	}
}

func TestPathsReturnsFreshSlice(t *testing.T) {
	tr := NewTree()
	tr.Add("a.txt")
	first := tr.Paths()
	first[0] = "mutated"
	second := tr.Paths()
	if second[0] != "a.txt" {
		t.Errorf("second Paths() = %v, want a.txt", second)
	}
}

func TestPathsEmptyTree(t *testing.T) {
	tr := NewTree()
	paths := tr.Paths()
	if len(paths) != 0 {
		t.Errorf("Paths() = %v, want empty", paths)
	}
}

func TestPathsSkipsDirectories(t *testing.T) {
	tr := NewTree()
	tr.Add("dir/file.txt")
	paths := tr.Paths()
	want := []string{"dir/file.txt"}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("Paths() = %v, want %v", paths, want)
	}
}
