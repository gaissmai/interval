// Package interval is an immutable datastructure for fast lookups in one dimensional intervals.
//
// The implementation is based on Treaps, augmented for interval.
//
// Immutability is achived because insert/upsert/delete will return a new Treap which will
// share some nodes with the original Treap.
// All nodes are read-only after creation, allowing concurrent readers to operate safely with concurrent writers.
//
// The time complexity is O(k*log(n)) where k is the number of returned items.
//
//  Insert()    O(log(n))
//  Upsert()    O(log(n))
//  Delete()    O(log(n))
//  Shortest()  O(log(n))
//  Largest()   O(log(n))
//
//  Subsets()   O(k*log(n))
//  Supersets() O(k*log(n))
//
// The space complexity is O(n).
//
// The author is propably the first (december 2022) using augmented Treaps
// as a very promising [data structure] for the representation of dynamic IP address tables
// for arbitrary ranges, that enables most and least specific range matching and even more lookup methods
// returning sets of intervals.
//
// The library can be used for all comparable one-dimensional intervals,
// but the author of the library uses it mainly for fast [IP range lookups] in access control lists (ACL)
// and in his own IP address management (IPAM) and network management software.
//
// The augmented Treap is NOT limited to IP CIDR ranges unlike the prefix trie.
// Arbitrary IP ranges and both IP versions can be processed together in this data structure.
//
// Due to the nature of Treaps the lookups and updates can be concurrently decoupled, without delayed rebalancing.
//
// To familiarize yourself with Treaps, see the extraordinarily good lectures from
// Pavel Mravin about Algorithms and Datastructures e.g. "[Treaps, implicit keys]"
// or follow [some links about Treaps] from one of the inventors.
//
// Especially useful is the paper "[Fast Set Operations Using Treaps]" by Guy E. Blelloch and Margaret Reid-Miller.
//
// [IP-Range lookups]: https://github.com/gaissmai/iprange
// [data structure]: https://ieeexplore.ieee.org/abstract/document/912716
// [Treaps, implicit keys]: https://youtu.be/svAHk-FAQgM
// [some links about Treaps]: http://faculty.washington.edu/aragon/treaps.html
// [Fast Set Operations Using Treaps]: https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf
package interval
