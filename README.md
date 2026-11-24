# go_path_tree

Constructs an in-memory tree representation of a directory hierarchy from a flat list of file paths.

```go
package main

import (
	"fmt"

	path_tree "go_path_tree"
)

func main() {
	tree := path_tree.NewTree()
	tree.AddAll([]string{
		"cmd/server/main.go",
		"internal/tree/tree.go",
		"README.md",
	})
	tree.Sort()

	for _, p := range tree.Paths() {
		fmt.Println(p)
	}
}
```

Output:

```
README.md
cmd/server/main.go
internal/tree/tree.go
```

## Why this exists

Given a flat list of paths (for example, the output of a file walker or a set of
generated file names), you often want to ask structural questions: which
directories exist, what is inside them, and what is the canonical list of files.
Building the tree once and walking it is simpler than re-scanning the list for
every query.

The trade-off is that the tree is purely additive and does not own the source
list. `Add` never removes a node, and adding a path that conflicts with an
existing file is a no-op past the conflict point rather than an error or a
replacement. That keeps the API small and the behaviour predictable; callers
that need validation can compare the input list against `Paths()`.

## API

- `type Node struct { Name string; IsDir bool; Children []*Node }`
- `type Tree struct { Root *Node }`
- `func NewTree() *Tree`
- `func (t *Tree) Add(path string)`
- `func (t *Tree) AddAll(paths []string)`
- `func (t *Tree) Sort()`
- `func (t *Tree) Paths() []string`

Paths are split on `/`. Empty components produced by leading, trailing, or
repeated slashes are ignored, so `a//b/` and `a/b` are the same path. A path
that resolves to no components is a no-op.

Intermediate components are always directories. The final component is always a
file. `Paths` returns only files, in depth-first pre-order, and `Sort` orders
children by name (files and directories interleaved).

## The awkward edge

If you add `a` and then `a/b`, the second call stops at `a`. The first call
already declared `a` to be a file, and this library will not silently convert it
into a directory and discard that fact. The result is a tree where `a` is a file
with no children, and `a/b` is not represented. If your input list can contain
both a file and a descendant of that file, decide which one wins before calling
`Add`.
