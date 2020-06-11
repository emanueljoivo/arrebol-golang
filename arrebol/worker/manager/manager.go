package manager

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"strconv"
)

type Manager struct {
	storage *storage.Storage
}

func NewManager(storage *storage.Storage) *Manager {
	return &Manager{
		storage:storage,
	}}

func (m *Manager) Join(w worker.Worker) (uint, error) {
	queueId := m.selectQueue(w)
	w.QueueID = queueId
	queue, err := m.storage.RetrieveQueue(queueId)
	if err != nil {
		return 0, err
	}
	queue.Workers = append(queue.Workers, &w)
	err = m.storage.SaveQueue(queue)
	if err != nil {
		return 0, err
	}
	log.Println("Worker [" + w.ID + "] has been assigned to queue [" + strconv.Itoa(int(queueId)) + "]")
	return queueId, nil
}

func (m *Manager) selectQueue(w worker.Worker) uint {
	log.Println("Selecting a queue for worker [" + w.ID + "]")
	return 1
}