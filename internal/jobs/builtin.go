package jobs

import (
	"fmt"
	"sync"

	"golang.org/x/term"
)

type Job struct {
	ID      int
	PID     int
	Command string
	State   string
}

type Manager struct {
	mu     sync.Mutex
	jobs   map[int]*Job
	nextId int
}

func NewManager() *Manager {
	return &Manager{
		jobs:   make(map[int]*Job),
		nextId: 1,
	}
}

func (m *Manager) Add(pid int, cmd string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextId
	m.jobs[id] = &Job{
		ID:      id,
		PID:     pid,
		Command: cmd,
		State:   "Running",
	}
	m.nextId++
	return id
}

func (m *Manager) MarkDone(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if job, exists := m.jobs[id]; exists {
		job.State = "Done"
	}
}

func (m *Manager) ListJobs(t *term.Terminal) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 1; i < m.nextId; i++ {
		if job, exists := m.jobs[i]; exists {
			fmt.Fprintf(t, "[%d] %s \t %s\n", job.ID, job.State, job.Command)

			if job.State == "Done" {
				delete(m.jobs, i)
			}
		}
	}
}
