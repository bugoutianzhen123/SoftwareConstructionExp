package cache

import (
    "sync"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
)

type MemoryCache struct {
    mu          sync.RWMutex
    userByID    map[int64]*domain.User
    users       []*domain.User
    projectByID map[int64]*domain.Project
    projects    []*domain.Project
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{userByID: map[int64]*domain.User{}, projectByID: map[int64]*domain.Project{}}
}

func (c *MemoryCache) GetUser(id int64) *domain.User {
    c.mu.RLock(); defer c.mu.RUnlock()
    return c.userByID[id]
}
func (c *MemoryCache) SetUser(u *domain.User) {
    c.mu.Lock(); defer c.mu.Unlock()
    c.userByID[u.ID] = u
}
func (c *MemoryCache) ListUsers() []*domain.User {
    c.mu.RLock(); defer c.mu.RUnlock()
    return c.users
}
func (c *MemoryCache) SetUsers(us []*domain.User) {
    c.mu.Lock(); defer c.mu.Unlock()
    c.users = us
}
func (c *MemoryCache) InvalidateUsers() {
    c.mu.Lock(); defer c.mu.Unlock()
    c.users = nil
}

func (c *MemoryCache) GetProject(id int64) *domain.Project {
    c.mu.RLock(); defer c.mu.RUnlock()
    return c.projectByID[id]
}
func (c *MemoryCache) SetProject(p *domain.Project) {
    c.mu.Lock(); defer c.mu.Unlock()
    c.projectByID[p.ID] = p
}
func (c *MemoryCache) ListProjects() []*domain.Project {
    c.mu.RLock(); defer c.mu.RUnlock()
    return c.projects
}
func (c *MemoryCache) SetProjects(ps []*domain.Project) {
    c.mu.Lock(); defer c.mu.Unlock()
    c.projects = ps
}
func (c *MemoryCache) InvalidateProjects() {
    c.mu.Lock(); defer c.mu.Unlock()
    c.projects = nil
}
