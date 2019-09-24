package beanstalkg

import (
	"fmt"
	"strings"
)

func (c *Connection) getYMLResp(cmd string) ([]byte, error) {
	resp, err := c.GetResp(cmd)
	if err != nil {
		return nil, err
	}

	var bodyLen int

	switch {
	case strings.HasPrefix(resp, "OK"):
		_, err = fmt.Sscanf(resp, "OK %d\r\n", &bodyLen)
		if err != nil {
			return nil, err
		}
	default:
		return nil, stringToError(resp)
	}

	return c.readBody(bodyLen)
}


func (c *Connection) StatsJob(id int) ([]byte, error) {
	cmd := fmt.Sprintf("stats-job %d\r\n", id)
	return c.getYMLResp(cmd)
}

func (c *Connection) StatsTube(tubename string) (*TubeStats, error) {
	cmd := fmt.Sprintf("stats-tube %s\r\n", tubename)
	b, err := c.getYMLResp(cmd)
	if err != nil {
		return nil, err
	}
	return tubeStatsParser(b), nil
}

func (c *Connection) Stats() (*InstanceStats, error) {
	b, err := c.getYMLResp("stats\r\n")
	if err != nil {
		return nil, err
	}
	res := statsParser(b)
	return res, nil
}

func (c *Connection) ListTubes() ([]string, error) {
	b, err := c.getYMLResp("list-tubes\r\n")
	if err != nil {
		return []string{}, err
	}
	return listParser(b), nil
}