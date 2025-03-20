package executor

import (
	"sync"
)

// Thay thế mutex toàn cục bằng một hệ thống quản lý boxID
var boxIDManager = NewBoxIDManager()

type BoxIDManager struct {
	mu      sync.Mutex
	usedIDs map[int]bool
}

func NewBoxIDManager() *BoxIDManager {
	return &BoxIDManager{
		usedIDs: make(map[int]bool),
	}
}

func (m *BoxIDManager) Acquire() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id := 1; id < 1000; id++ {
		if !m.usedIDs[id] {
			m.usedIDs[id] = true
			return id
		}
	}
	return -1
}

func (m *BoxIDManager) Release(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.usedIDs, id)
}
