
# bursavich.dev/arc

[![License](https://img.shields.io/badge/license-mit-blue.svg?style=for-the-badge)](https://raw.githubusercontent.com/abursavich/arc/main/LICENSE)
[![GoDev Reference](https://img.shields.io/static/v1?logo=go&logoColor=white&color=00ADD8&label=dev&message=reference&style=for-the-badge)](https://pkg.go.dev/bursavich.dev/arc)
[![Go Report Card](https://goreportcard.com/badge/bursavich.dev/arc?style=for-the-badge)](https://goreportcard.com/report/bursavich.dev/arc)

Package arc implements a generic adaptive replacement cache.

It's similar to IBM's patented version in the sense that it pivots (or adapts) between a most-recently-used
and frequently-used cache when misses hit a ghost cache of recently evicted entries. However, this version
permits the deletion of entries and the precise details of pivoting rate and ghost cache eviction are different.
Given the same sequence of mutually-supported operations, the contents of the two implementations may diverge.

I am not a lawyer. This is not legal advice.

See:
- [Adaptive Replacement Cache](https://en.wikipedia.org/wiki/Adaptive_replacement_cache)
- [ARC: A Self-tuning, Low Overhead Replacement Cache](https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf)
- [Patent US6996676 - System and method for implementing an adaptive replacement cache policy](https://www.google.com/patents/US6996676)
