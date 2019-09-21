package beanstalkg

import (
	"fmt"
	"strings"
	"time"
)

func assertExpected(expected, got string) error {
	if expected != got {
		return UnexpectedResponse
	}
	return nil
}

func handleReserve(c *Connection, cmd string) (*Job, error) {
	resp, err := c.GetResp(cmd)
	if err != nil {
		return nil, err
	}
	var id int
	var bodyLen int
	if strings.HasPrefix(resp, "RESERVED") {
		_, err := fmt.Sscanf(resp, "RESERVED %d %d\r\n", &id, &bodyLen)
		fmt.Println(id, bodyLen)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, stringToError(resp)
	}
	body, err := c.readBody(bodyLen)
	return &Job{ID: id, Body: body, conn: c}, err
}

func handlePut(c *Connection, cmd string) (int, error) {
	resp, err := c.GetResp(cmd)
	if err != nil {
		return 0, err
	}
	switch {
	case strings.HasPrefix(resp, "INSERTED"):
		var id int
		_, err = fmt.Sscanf(resp, "INSERTED %d\r\n", &id)
		return id, err
	case strings.HasPrefix(resp, "BURIED"):
		var id int
		_, err := fmt.Sscanf(resp, "BURIED %d\r\n", &id)
		if err != nil {
			return 0, UnexpectedResponse
		}
		return id, Buried
	default:
		return 0, stringToError(resp)
	}
}

func (c *Connection) ReserveWithTimeout(timeout time.Duration) (*Job, error) {
	cmd := fmt.Sprintf("reserve-with-timeout %d\r\n", int(timeout.Seconds()))
	return handleReserve(c, cmd)
}

func (c *Connection) Reserve() (*Job, error) {
	return handleReserve(c, "reserve\r\n")
}

func (c *Connection) Ignore(tubename string) (int, error) {
	cmd := fmt.Sprintf("ignore %s\r\n", tubename)
	resp, err := c.GetResp(cmd)
	if err != nil {
		return -1, err
	}
	var tubeNum int
	_, err = fmt.Sscanf(resp, "WATCHING %d\r\n", &tubeNum)
	if err != nil {
		return -1, stringToError(resp)
	}
	return tubeNum, nil
}

func (c *Connection) Watch(tubename string) (int, error) {
	cmd := fmt.Sprintf("watch %s\r\n", tubename)

	resp, err := c.GetResp(cmd)
	if err != nil {
		return -1, err
	}
	var tubeNum int
	_, err = fmt.Sscanf(resp, "WATCHING %d\r\n", &tubeNum)
	if err != nil {
		return -1, stringToError(tubename)
	}
	return tubeNum, nil
}

func (c *Connection) Use(tubename string) error {
	if len(tubename) > 200 {
		return InvalidTubeName
	}
	cmd := fmt.Sprintf("use %s", tubename)
	res, _ := c.GetResp(cmd)
	fmt.Println(res)
	return assertExpected(fmt.Sprintf("USING %s\r\n", tubename), res)
}

func (c *Connection) PutBytes(data []byte, priority int, delay, timeToRun time.Duration) (int, error) {
	cmd := fmt.Sprintf("put %d %d %d %d\r\n%s\r\n", priority, uint64(delay.Seconds()), uint64(timeToRun.Seconds()), len(data), string(data))
	return handlePut(c, cmd)
}

func (c *Connection) PutString(data string, priority int, delay, timeToRun time.Duration) (int, error) {
	cmd := fmt.Sprintf("put %d %d %d %d\r\n%s\r\n", priority, uint64(delay.Seconds()), uint64(timeToRun.Seconds()), len(data), data)
	return handlePut(c, cmd)
}

func (c *Connection) Release(id int, priority int, delay time.Duration) error {
	cmd := fmt.Sprintf("release %d %d %d\r\n", id, priority, int(delay.Seconds()))
	resp, err := c.GetResp(cmd)
	if err != nil {
		return err
	}
	return assertExpected("RELEASED\r\n", resp)
}

func (c *Connection) Bury(id, priority int) error {
	cmd := fmt.Sprintf("bury %d %d\r\n", id, priority)
	resp, err := c.GetResp(cmd)
	if err != nil {
		return err
	}
	return assertExpected("BURIED\r\n", resp)
}

func (c *Connection) KickJob(id int) (error) {
	cmd := fmt.Sprintf("kick-job %d\r\n", id)
	resp, err := c.GetResp(cmd)
	if err != nil {
		return err
	}
	return assertExpected("KICKED\r\n", resp)
}

func (c *Connection) Kick(maxJobs int) (int, error) {
	cmd := fmt.Sprintf("kick %d\r\n", maxJobs)

	resp, err := c.GetResp(cmd)
	if err != nil {
		return 0, err
	}

	var id int
	if strings.HasPrefix(resp, "KICKED") {
		_, err := fmt.Sscanf(resp, "KICKED %d\r\n", &id)
		if err != nil {
			return 0, UnexpectedResponse
		}
		return id, nil
	}
	return 0, stringToError(resp)
}

func (c *Connection) Touch(id int) error {
	cmd := fmt.Sprintf("touch %d\r\n", id)
	resp, err := c.GetResp(cmd)
	if err != nil {
		return err
	}
	return assertExpected("TOUCHED\r\n", resp)
}

func (c *Connection) Delete(id int) error {
	cmd := fmt.Sprintf("delete %d\r\n", id)
	resp, err := c.GetResp(cmd)
	if err != nil {
		return err
	}
	return assertExpected("DELETED\r\n", resp)
}

func (c *Connection) Quit() {
	defer c.connection.Close()
	_, _ = c.GetResp("quit \r\n")
}

