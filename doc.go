// Package linq provides lazy, statically typed LINQ-style queries for Go 1.27.
//
// Generic methods allow each fluent operation to infer its own result, key,
// inner, or accumulator types without reflection or type assertions. Queries
// can be created from slices, maps, channels, strings, iter.Seq values, and
// custom Iterable collections.
package linq
