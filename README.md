# Generic Data Structures

With the release of Go 1.18, it is possible to implement generic data
structures in Go. This repository contains some data structures I have found
useful. See the individual directories for more information about each data
structure.

* `avl`: an AVL tree.
* `cache`: a wrapper around `map[K]V` that uses a maximum size and evicts
  elements using LRU when full.
* `hashmap`: a hashmap with linear probing. The main feature is that
  the hashmap can be efficiently copied, using copy-on-write under the hood.
* `list`: a doubly-linked list.
* `rope`: a generic rope, which is similar to an array but supports efficient
  insertion and deletion from anywhere in the array. Ropes are typically used
  for arrays of bytes, but this rope is generic.
* `stack`: a LIFO stack.

This project is currently in-progress.

Planned additions:

* Interval tree
* Hashset
* Trie
* Queue
* Approximate membership query filters
* B-tree
* Better tests
* Examples
* Benchmarks
