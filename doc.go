// Package interval is an immutable datastructure for fast lookups in one dimensional intervals.
//
// The implementation is based on treaps, augmented for intervals. Treaps are randomized self balancing binary search trees.
//
// Immutability is achieved because insert/upsert/delete will return a new treap which will
// share some nodes with the original treap. All nodes are read-only after creation,
// allowing concurrent readers to operate safely with concurrent writers.
//
// The time complexity is O(k*log(n)) where k (k <= n) is the number of returned items.
//
//  Insert()    O(log(n))
//  Delete()    O(log(n))
//  Shortest()  O(log(n))
//  Largest()   O(log(n))
//
//  Subsets()   O(k*log(n))
//  Supersets() O(k*log(n))
//
// The space complexity is O(n).
//
// The library can be used for all comparable one-dimensional intervals,
// but the author of the library uses it mainly for fast IP range lookups in access control lists (ACL)
// and in his own IP address management (IPAM) and network management software.
//
// The augmented treap is NOT LIMITED to IP CIDR ranges unlike the prefix trie.
// Arbitrary IP ranges and both IP versions (v4/v6) can be processed transparently together in the same treap.
//
// Due to the nature of treaps the lookups and updates can be concurrently decoupled,
// without delayed rebalancing, promising to be a perfect match for a software-router or firewall.
package interval
