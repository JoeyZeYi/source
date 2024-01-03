package util

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: sync.Map{},
	}
}

func (m *Map[K, V]) Get(key K) *V {
	value, ok := m.m.Load(key)
	if ok {
		if val, valOk := value.(*V); valOk {
			return val
		}
	}
	return nil
}

func (m *Map[K, V]) Delete(k K) {
	m.m.Delete(k)
}

func (m *Map[K, V]) Set(k K, t *V) {
	m.m.Store(k, t)
}

func (m *Map[K, V]) CompareAndSwap(k K, old, new *V) bool {
	return m.m.CompareAndSwap(k, old, new)
}

func (m *Map[K, V]) CompareAndDelete(k K, old *V) bool {
	return m.m.CompareAndDelete(k, old)
}

func (m *Map[K, V]) GetAndDelete(k K) (*V, bool) {
	value, loaded := m.m.LoadAndDelete(k)
	if val, valOk := value.(*V); valOk {
		return val, loaded
	}
	return nil, loaded
}

func (m *Map[K, V]) GetOrSet(k K, v *V) (*V, bool) {
	value, loaded := m.m.LoadOrStore(k, v)
	if val, valOk := value.(*V); valOk {
		return val, loaded
	}
	return nil, loaded
}

func (m *Map[K, V]) Swap(k K, v *V) (*V, bool) {
	value, loaded := m.m.Swap(k, v)
	if val, valOk := value.(*V); valOk {
		return val, loaded
	}
	return nil, loaded
}

func (m *Map[K, V]) Maps() map[K]*V {
	result := make(map[K]*V)
	m.m.Range(func(key, value any) bool {
		if val, valOk := value.(*V); valOk {
			k, kOk := key.(K)
			if kOk {
				result[k] = val
			}
		}
		return true
	})
	return result
}

func (m *Map[K, V]) Range(f func(key, value any) bool) {
	m.m.Range(f)
}
