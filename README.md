# Generic Data Structures

With the release of Go 1.18, it will be possible to implement generic data
structures in Go. This repository contains some data structures I have found
useful implemented with generics. See the individual directories for more
information about each data structure.

* `avl`: an AVL tree.
* `btree`: a B-tree.
* `cache`: a wrapper around `map[K]V` that uses a maximum size and evicts
  elements using LRU when full.
* `hashmap`: a hashmap with linear probing. The main feature is that
  the hashmap can be efficiently copied, using copy-on-write under the hood.
* `hashset`: a hashset that uses the hashmap as the underlying storage.
* `interval`: an interval tree, implemented as an augmented AVL tree.
* `list`: a doubly-linked list.
* `rope`: a generic rope, which is similar to an array but supports efficient
  insertion and deletion from anywhere in the array. Ropes are typically used
  for arrays of bytes, but this rope is generic.
* `stack`: a LIFO stack.
* `trie`: a ternary search trie.

The package also includes support for iterators, in the `iter` subpackage.
Most data structures provide an iterator API, which can be used with some
convenience functions in `iter`.

See each subpackage for documentation and examples. The top-level `generic`
package provides some useful types and constraints. See [DOC.md](DOC.md)
for documentation.

This project is currently in-progress and the API is not stable. A stable
version will be released when Go 1.18 is released. Currently, this serves
as an experiment to get familiar with Go's generics, and figure out what
approaches for using them are best.

# Discussion

We are in the very early stages of generics in Go and it is not clear what the
best practices are. This project is an attempt to become familiar with Go's
generics and determine what works well and what doesn't. If you have feedback
on the implementation, please open an issue for further discussion.

Some notes:

* Iterators: the `iter` package provides an API for data structures to return
  iterators. The main use of this is to loop over all elements in a data
  structure and apply some function. Is this better than just returning all
  key-value pairs as a slice, or using a `.Each(callback)` function? In some
  cases, the iterator can be lazy, which is better than returning a slice
  because the slice has to be pre-computed. But for other data structures
  (especially trees), it is difficult to make a lazy iterator, and the iterator
  doesn't provide much benefit (has to either pre-allocate a slice of results,
  or allocate many function closures). The `.Each` approach doesn't ever
  allocate anything but is less flexible than the other approaches. At the
  moment I am not sure iterators are worth having.

* Custom key types: since there are no methods defined on primitive types, and
  no operators defined on user types, it is impossible to make a generic data
  structures/function with operator constraints that accepts both primitive and
  user types. Should we use the primitive type constraints as much as possible,
  or is having wrapper types for all primitive types like in this package the
  correct approach? For example, the AVL tree constrains keys to be
  `generic.Lesser[T]`, but they could be `constraints.ordered`. This would mean
  keys must be primitive types. As it currently stands, to use a primitive type
  as a key with `generic.Lesser[T]`, one must use the appropriate wrapper from
  the `generic` package: `generic.Int(x)`. Note: in some cases, we must always
  use a wrapper, such as for data structures that require hashable keys.

# Contributing

There are more data structures that may be useful to have, such as bloom
filters, queues, graph representations, and more kinds of search trees.
If you would like to contribute a data structure please let me know.

It would also be useful to have comprehensive benchmarks for the data
structures, comparing to standard library implementations when possible, or
just on their own. Benchmarks will also allow us to profile and optimize
the implementations.
