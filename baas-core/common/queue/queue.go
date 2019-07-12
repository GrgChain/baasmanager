package queue

type Item interface {
}

type ItemQueue interface {
	Init() Queue
	Enqueue(t Item)
	Dequeue() *Item
	IsEmpty() bool
	Size() int
}

// Item the type of the queue
type Queue struct {
	items []Item
}

// New creates a new ItemQueue
func (s *Queue) Init() *Queue {
	s.items = []Item{}
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *Queue) Enqueue(t Item) {
	s.items = append(s.items, t)
}

// dequeue
func (s *Queue) Dequeue() Item {
	item := s.items[0] // 先进先出
	s.items = s.items[1:len(s.items)]
	return item
}

func (s *Queue) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of Items in the queue
func (s *Queue) Size() int {
	return len(s.items)
}

func NewQueue() *Queue {
	que := &Queue{}
	que.Init()
	return que
}
