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
* [`interval`](./interval): an interval tree, implemented as an augmented AVL tree.
* [`list`](./list): a doubly-linked list.
* [`rope`](./rope): a generic rope, which is similar to an array but supports efficient
  insertion and deletion from anywhere in the array. Ropes are typically used
  for arrays of bytes, but this rope is generic.
* [`stack`](./stack): a LIFO stack.
* [`trie`](./trie): a ternary search trie.
* [`queue`](./queue): a First In First Out (FIFO) queue.

See each subpackage for documentation and examples. The top-level `generic`
package provides some useful types and constraints. See [DOC.md](DOC.md) for
documentation.

# Contributing

There are more data structures that may be useful to have, such as bloom
filters, graph representations, and more kinds of search trees.
If you would like to contribute a data structure please let me know.

It would also be useful to have comprehensive benchmarks for the data
structures, comparing to standard library implementations when possible, or
just on their own. Benchmarks will also allow us to profile and optimize
the implementations.
