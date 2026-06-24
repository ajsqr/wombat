package cmap

import "sync"

type ConcurentMap[k comparable, v any] interface {
	Get(key k) (v, bool)
	Set(key k, value v)
	Delete(key k)
}

func New[k comparable, v any]() *Cmap[k, v] {
	data := make(map[k]v)
	return &Cmap[k, v]{
		data: data,
		mx:   &sync.RWMutex{},
	}

}

type Cmap[k comparable, v any] struct {
	data map[k]v
	mx   *sync.RWMutex
}

func (cm *Cmap[k, v]) Get(key k) (v, bool) {
	cm.mx.RLock()
	defer cm.mx.RUnlock()
	value, ok := cm.data[key]
	return value, ok
}

func (cm *Cmap[k, v]) Set(key k, value v) {
	cm.mx.Lock()
	defer cm.mx.Unlock()
	cm.data[key] = value
}

func (cm *Cmap[k, v]) Delete(key k) {
	cm.mx.Lock()
	defer cm.mx.Unlock()
	delete(cm.data, key)
}
