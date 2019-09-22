package beanstalkg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InstanceStats struct {
	RusageUtime float64
	Uptime int
	BinlogCurrentIndex int
	BinlogRecordsMigrated int
	CmdStatsJob int
	CmdListTubes int
	CmdStats int
	CmdPeekReady int
	CmdDelete int
	CmdReserveWithTimeout int
	JobTimeouts int
	MaxJobSize int
	RusageStime float64
	Id string
	CurrentJobsBuried int
	CmdPeekDelayed int
	CmdTouch int
	CurrentConnections int
	BinlogRecordsWritten int
	CurrentJobsReady int
	CmdReserve int
	CmdPeek int
	CmdRelease int
	CmdBury int
	TotalJobs int
	Version int
	CurrentJobsUrgent int
	CmdPut int
	CmdStatsTube int
	CmdListTubesWatched int
	CurrentWorkers int
	BinlogOldestIndex int
	BinlogMaxSize int
	Hostname string
	CurrentJobsReserved int
	CmdWatch int
	CmdIgnore int
	TotalConnections int
	CurrentJobsDelayed int
	CmdPeekBuried int
	CmdListTubeUsed int
	CmdPauseTube int
	CurrentTubes int
	CurrentProducers int
	CurrentWaiting int
	Pid int
	CmdUse int
	CmdKick int
}

func (s *InstanceStats) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		key := normaliseKey(k)
		_ = SetField(s, key, v)
	}
	return nil
}

func normaliseKey(key string) string {
	noHashes := strings.Replace(key, "-", " ", -1)
	titleCase := strings.Title(noHashes)
	return strings.Replace(titleCase, " ", "", -1)
}

var stringFields = map[string]bool{
	"id": true,
	"hostname": true,
}

var floatFields = map[string]bool {
	"rusage-utime": true,
	"rusage-stime": true,
}

func isFieldType(key string, fields map[string]bool) bool {
	_, ok := fields[key]
	return ok
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func statsParser(res []byte) *InstanceStats {
	instanceStats := make(map[string]interface{})
	resultString := string(res)
	results := strings.Split(resultString, "\n")
	for _, r := range results {
		stat := strings.Split(r, ":")
		if len(stat) == 2 {
			key := stat[0]
			value := strings.TrimSpace(stat[1])
			switch {
			case isFieldType(key, stringFields):
				instanceStats[key] = value
			case isFieldType(key, floatFields):
				floatValue, _ := strconv.ParseFloat(value, 64)
				instanceStats[key] = floatValue
			default:
				intValue, _ := strconv.ParseInt(value, 10, 64)
				instanceStats[key] = intValue
			}
		}
	}
	result := &InstanceStats{}
	result.FillStruct(instanceStats)
	return result
}
