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
	mu   sync.Mutex
	jobs map[int]*Job
}

func NewManager() *Manager {
	return &Manager{
		jobs: make(map[int]*Job),
	}
}

func (m *Manager) Add(pid int, cmd string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := 1
	for {
		if _, exists := m.jobs[id]; !exists {
			break
		}
		id++
	}
	m.jobs[id] = &Job{
		ID:      id,
		PID:     pid,
		Command: cmd,
		State:   "Running                 ",
	}
	return id
}

func (m *Manager) MarkDone(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if job, exists := m.jobs[id]; exists {
		job.State = "Done                    "
	}
}

func (m *Manager) ReapJobs(t *term.Terminal) {
	m.mu.Lock()
	defer m.mu.Unlock()

	maxId := 0
	for id := range m.jobs {
		if id > maxId {
			maxId = id
		}
	}
	var curId, prevId int
	for id := 1; id < maxId; id++ {
		if _, exists := m.jobs[id]; exists {
			prevId = curId
			curId = id
		}
	}
	for i := 1; i < maxId; i++ {
		if job, exists := m.jobs[i]; exists {
			if job.State == "Done                    " {
				var sym string
				switch job.ID {
				case curId:
					sym = "+"
				case prevId:
					sym = "-"
				default:
					sym = " "
				}
				var cmd = job.Command
				if job.State == "Done                    " {
					cmd = cmd[:len(cmd)-1]
				}
				fmt.Fprintf(t, "[%d]%s  %s%s\n", job.ID, sym, job.State, cmd)

				delete(m.jobs, i)

			}
		}
	}
}

func (m *Manager) ListJobs(t *term.Terminal) {
	m.mu.Lock()
	defer m.mu.Unlock()

	maxId := 0
	for id := range m.jobs {
		if id > maxId {
			maxId = id
		}
	}

	var curId, prevId int
	for id := 1; id < maxId; id++ {
		if _, exists := m.jobs[id]; exists {
			prevId = curId
			curId = id
		}
	}
	for i := 1; i < maxId; i++ {
		if job, exists := m.jobs[i]; exists {
			var sym string
			switch job.ID {
			case curId:
				sym = "+"
			case prevId:
				sym = "-"
			default:
				sym = " "
			}
			var cmd = job.Command
			if job.State == "Done                    " {
				cmd = cmd[:len(cmd)-1]
			}
			fmt.Fprintf(t, "[%d]%s  %s%s\n", job.ID, sym, job.State, cmd)

			// if job.State == "Done                    " {
			// 	delete(m.jobs, i)
			// }
		}
	}
}
