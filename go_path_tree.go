// Package go_path_tree constructs an in-memory tree representation of a
// directory hierarchy from a flat list of file paths.
package go_path_tree

import (
	"sort"
	"strings"
)

// Node is a single entry in the path tree. A Node may represent either a
// directory or a file. The distinction is made by the IsDir field.
type Node struct {
	// Name is the final path component of this node (the base name).
	Name string
	// IsDir reports whether this node represents a directory.
	IsDir bool
	// Children holds the child nodes of a directory. It is nil for files.
	Children []*Node
}

// Tree is an in-memory representation of a directory hierarchy.
type Tree struct {
	// Root is the synthetic root node. It is always a directory and its Name
	// is empty.
	Root *Node
}

// NewTree returns an empty Tree with a synthetic directory root.
func NewTree() *Tree {
	return &Tree{Root: &Node{Name: "", IsDir: true}}
}

// Add inserts a single path into the tree. Paths are split on '/'. Empty
// components (produced by leading, trailing or repeated slashes) are ignored,
// so "a//b/" and "a/b" describe the same entry. A path that resolves to no
// components is a no-op.
//
// Intermediate components are always created as directories. The final
// component is created as a file. If a component already exists as a file,
// Add stops there rather than trying to descend into it; this keeps the tree
// internally consistent when a list mixes "a" and "a/b".
func (t *Tree) Add(path string) {
	parts := splitPath(path)
	if len(parts) == 0 {
		return
	}

	cur := t.Root
	for i, part := range parts {
		last := i == len(parts)-1
		child := findChild(cur, part)
		if child == nil {
			child = &Node{Name: part, IsDir: !last}
			cur.Children = append(cur.Children, child)
		}
		if !last && !child.IsDir {
			// A previous Add created this component as a file. We cannot
			// descend into a file, so stop. The alternative (turning the
			// file into a directory) would silently discard the earlier
			// entry and is not what callers expect from an additive API.
			return
		}
		cur = child
	}
}

// AddAll inserts every path in paths, in order. It is equivalent to calling
// Add once per element.
func (t *Tree) AddAll(paths []string) {
	for _, p := range paths {
		t.Add(p)
	}
}

// Sort recursively sorts every directory's children by name. Files and
// directories are interleaved rather than grouped; callers that need a
// particular presentation can walk the tree and regroup.
func (t *Tree) Sort() {
	sortNode(t.Root)
}

// Paths returns every file path in the tree, in depth-first pre-order. Only
// nodes with IsDir == false are emitted; directories are not listed. The
// returned slice is fresh and safe to modify.
func (t *Tree) Paths() []string {
	var out []string
	collectPaths(t.Root, "", &out)
	return out
}

func splitPath(path string) []string {
	raw := strings.Split(path, "/")
	parts := raw[:0]
	for _, p := range raw {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func findChild(n *Node, name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func sortNode(n *Node) {
	if len(n.Children) == 0 {
		return
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
	for _, c := range n.Children {
		if c.IsDir {
			sortNode(c)
		}
	}
}

func collectPaths(n *Node, prefix string, out *[]string) {
	for _, c := range n.Children {
		full := c.Name
		if prefix != "" {
			full = prefix + "/" + c.Name
		}
		if c.IsDir {
			collectPaths(c, full, out)
		} else {
			*out = append(*out, full)
		}
	}
}
