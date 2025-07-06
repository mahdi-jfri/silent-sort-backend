package hub

import "sync"

type HubKeeper struct {
	mutex   sync.RWMutex
	hubById map[string]*Hub
}

func NewHubKeeper() *HubKeeper {
	return &HubKeeper{
		hubById: make(map[string]*Hub),
	}
}

func (hk *HubKeeper) GetHub(hubId string) *Hub {
	hk.mutex.RLock()
	defer hk.mutex.RUnlock()
	return hk.hubById[hubId]
}

func (hk *HubKeeper) SetHub(hubId string, hub *Hub) {
	hk.mutex.Lock()
	defer hk.mutex.Unlock()
	hk.hubById[hubId] = hub
	return
}
