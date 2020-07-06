package storage

func (s *Storage) SaveQueue(q *Queue) error {
	return s.driver.Save(&q).Error
}

func (s *Storage) RetrieveQueue(queueID uint) (*Queue, error) {
	var queue Queue
	err := s.driver.First(&queue, queueID).Error
	queue.Jobs, _ = s.RetrieveJobsByQueueID(queueID)
	queue.Workers, _ = s.RetrieveWorkersByQueueID(queueID)
	// TODO Retrieve resource nodes of queue
	return &queue, err
}

func (s *Storage) RetrieveQueues() ([]*Queue, error) {
	var queues []*Queue

	err := s.driver.Find(&queues).Error

	return queues, err
}

func (s *Storage) GetDefaultQueue() (*Queue, error) {
	var queue Queue
	const QIDDefault = 1
	if err := s.driver.Where("id = ?", QIDDefault).First(&queue).Error; err == nil {
		s.driver.First(&queue, 1)
	} else {
		return nil, err
	}
	return &queue, nil
}
