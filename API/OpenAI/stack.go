package openai

import "sync"

type Queue struct {
	Mtx   sync.Mutex
	Items []SavePrompt
}

func NewQueue() *Queue {
	return &Queue{
		Items: make([]SavePrompt, 0),
	}
}

func (q *Queue) Clear() {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	q.Items = []SavePrompt{}
}

func (q *Queue) Enqueue(i SavePrompt) {
	q.Mtx.Lock()
	q.Items = append(q.Items, i)
	q.Mtx.Unlock()
}

func (q *Queue) Dequeue() (bool, SavePrompt) {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	if len(q.Items) == 0 {
		return false, SavePrompt{}
	}
	var result = q.Items[0]
	q.Items = q.Items[1:]
	return true, result
}
