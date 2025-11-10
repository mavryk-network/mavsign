// Package hashmap implements a generic wrapper around a built in map type which allows
// interface types like PublicKeyHash to be used as map keys
package hashmap

import mv "github.com/mavryk-network/mavbingo/v2"

type KV[K, V any] struct {
	Key K
	Val V
}

// HashMap is a wrapper around a built in map type which allows interface types
// like mavbingo.PublicKeyHash to be used as map keys
type HashMap[H mv.Comparable[K], K mv.ToComparable[H, K], V any] map[H]V

func (m HashMap[H, K, V]) Insert(key K, val V) (V, bool) {
	k := key.ToComparable()
	v, ok := m[k]
	m[k] = val
	return v, ok
}

func (m HashMap[H, K, V]) Get(key K) (V, bool) {
	k := key.ToComparable()
	v, ok := m[k]
	return v, ok
}

func (m HashMap[H, K, V]) ForEach(cb func(key K, val V) bool) bool {
	for k, v := range m {
		if !cb(k.ToKey(), v) {
			return false
		}
	}
	return true
}

func New[H mv.Comparable[K], K mv.ToComparable[H, K], V any](init []KV[K, V]) HashMap[H, K, V] {
	m := make(HashMap[H, K, V])
	for _, kv := range init {
		m.Insert(kv.Key, kv.Val)
	}
	return m
}

type PublicKeyKV[V any] KV[mv.PublicKeyHash, V]

// PublicKeyHashMap is a shortcut for a map with mavbingo.PublicKeyHash keys
type PublicKeyHashMap[V any] HashMap[mv.EncodedPublicKeyHash, mv.PublicKeyHash, V]

func NewPublicKeyHashMap[V any](init []PublicKeyKV[V]) PublicKeyHashMap[V] {
	tmp := make([]KV[mv.PublicKeyHash, V], len(init))
	for i, x := range init {
		tmp[i] = KV[mv.PublicKeyHash, V](x)
	}
	return PublicKeyHashMap[V](New[mv.EncodedPublicKeyHash](tmp))
}

func (m PublicKeyHashMap[V]) Insert(key mv.PublicKeyHash, val V) (V, bool) {
	return HashMap[mv.EncodedPublicKeyHash, mv.PublicKeyHash, V](m).Insert(key, val)
}

func (m PublicKeyHashMap[V]) Get(key mv.PublicKeyHash) (V, bool) {
	return HashMap[mv.EncodedPublicKeyHash, mv.PublicKeyHash, V](m).Get(key)
}

func (m PublicKeyHashMap[V]) ForEach(cb func(key mv.PublicKeyHash, val V) bool) bool {
	return HashMap[mv.EncodedPublicKeyHash, mv.PublicKeyHash, V](m).ForEach(cb)
}
