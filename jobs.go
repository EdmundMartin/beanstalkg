package beanstalkg

import "time"

type Job struct {
	ID int
	Body []byte
	conn *Connection
}


func (j *Job) Release(priority int, delay time.Duration) error {
	return j.conn.Release(j.ID, priority, delay)
}

func (j *Job) Bury(priority int) error {
	return j.conn.Bury(j.ID, priority)
}

func (j *Job) KickJob() error {
	return j.conn.KickJob(j.ID)
}

func (j *Job) Touch() error {
	return j.conn.Touch(j.ID)
}

func (j *Job) Delete() error {
	return j.conn.Delete(j.ID)
}

func (j *Job) Stats() ([]byte, error) {
	return j.conn.StatsJob(j.ID)
}