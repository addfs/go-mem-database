package tools

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketNumber int) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, ticketNumber),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}
