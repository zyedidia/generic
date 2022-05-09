# Generic Data Structures

[![Go Reference](https://pkg.go.dev/badge/github.com/zyedidia/generic.svg)](https://pkg.go.dev/github.com/zyedidia/generic)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zyedidia/generic/blob/master/LICENSE)

This package implements some generic data structures.

* [`avl`](./avl): an AVL tree.
* [`btree`](./btree): a B-tree.
* [`cache`](./cache): a wrapper around `map[K]V` that uses a maximum size and evicts
  elements using LRU when full.
* [`hashmap`](./hashmap): a hashmap with linear probing. The main feature is that
  the hashmap can be efficiently copied, using copy-on-write under the hood.
* [`hashset`](./hashset): a hashset that uses the hashmap as the underlying storage.
* [`mapset`](./mapset): a set that uses Go's built-in map as the underlying storage.
* [`multimap`](./multimap): an associative container that permits multiple entries with the same key.
* [`interval`](./interval): an interval tree, implemented as an augmented AVL tree.
* [`list`](./list): a doubly-linked list.
* [`rope`](./rope): a generic rope, which is similar to an array but supports efficient
  insertion and deletion from anywhere in the array. Ropes are typically used
  for arrays of bytes, but this rope is generic.
* [`stack`](./stack): a LIFO stack.
* [`trie`](./trie): a ternary search trie.
* [`queue`](./queue): a First In First Out (FIFO) queue.
* [`heap`](./heap): a binary heap.

See each subpackage for documentation and examples. The top-level `generic`
package provides some useful types and constraints. See [DOC.md](DOC.md) for
documentation.

# Contributing

If you would like to contribute a new feature, please let me know first what
you would like to add (via email or issue tracker). Here are some ideas:

* New data structures (bloom filters, graph structures, concurrent data
  structures, adaptive radix tree, or other kinds of search trees).
* Benchmarks, and optimization of the existing data structures based on those
  benchmarks. The hashmap is an especially good target.
* Design and implement a nice iterator API.
* Improving tests (perhaps we can use Go's new fuzzing capabilities).
