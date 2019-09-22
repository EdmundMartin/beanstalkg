package beanstalkg

import "time"

type Job struct {
	ID int
	Body []byte
	conn *Connection
}

// Release will release the given job
func (j *Job) Release(priority int, delay time.Duration) error {
	return j.conn.Release(j.ID, priority, delay)
}

// Bury will bury the given job
func (j *Job) Bury(priority int) error {
	return j.conn.Bury(j.ID, priority)
}

// KickJob will kick the given job
func (j *Job) KickJob() error {
	return j.conn.KickJob(j.ID)
}

// Touch will touch the given job - giving more time to work
func (j *Job) Touch() error {
	return j.conn.Touch(j.ID)
}

// Delete will delete the given job
func (j *Job) Delete() error {
	return j.conn.Delete(j.ID)
}

func (j *Job) Stats() ([]byte, error) {
	return j.conn.StatsJob(j.ID)
}