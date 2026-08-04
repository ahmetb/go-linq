// Package linq provides methods for querying and manipulating slices, arrays,
// maps, strings, channels and collections.
//
// Starting with v5, go-linq is fully type-safe: Query[T] is a generic type and
// type-changing operators such as Select, Join and GroupBy are generic methods
// (a Go 1.27 language feature), so queries compose without interface{} values,
// type assertions or reflection.
//
// Authors: Alexander Kalankhodzhaev (kalan), Ahmet Alp Balkan, Cleiton Marques
// Souza.
package linq
