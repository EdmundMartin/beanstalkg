package beanstalkg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InstanceStats struct {
	RusageUtime           float64
	Uptime                int64
	BinlogCurrentIndex    int64
	BinlogRecordsMigrated int64
	CmdStatsJob           int64
	CmdListTubes          int64
	CmdStats              int64
	CmdPeekReady          int64
	CmdDelete             int64
	CmdReserveWithTimeout int64
	JobTimeouts           int64
	MaxJobSize            int64
	RusageStime           float64
	Id                    string
	CurrentJobsBuried     int64
	CmdPeekDelayed        int64
	CmdTouch              int64
	CurrentConnections    int64
	BinlogRecordsWritten  int64
	CurrentJobsReady      int64
	CmdReserve            int64
	CmdPeek               int64
	CmdRelease            int64
	CmdBury               int64
	TotalJobs             int64
	Version               int64
	CurrentJobsUrgent     int64
	CmdPut                int64
	CmdStatsTube          int64
	CmdListTubesWatched   int64
	CurrentWorkers        int64
	BinlogOldestIndex     int64
	BinlogMaxSize         int64
	Hostname              string
	CurrentJobsReserved   int64
	CmdWatch              int64
	CmdIgnore             int64
	TotalConnections      int64
	CurrentJobsDelayed    int64
	CmdPeekBuried         int64
	CmdListTubeUsed       int64
	CmdPauseTube          int64
	CurrentTubes          int64
	CurrentProducers      int64
	CurrentWaiting        int64
	Pid                   int64
	CmdUse                int64
	CmdKick               int64
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
	"id":       true,
	"hostname": true,
}

var floatFields = map[string]bool{
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
	results := strings.Split(string(res), "\n")
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

func listParser(res []byte) []string {
	tubes := []string{}
	results := strings.Split(string(res), "\n")
	for _, res := range results {
		if len(res) > 2 && res[:2] == "- " {
			tubes = append(tubes, strings.TrimSpace(res[2:]))
		}
	}
	return tubes
}

type TubeStats struct {
	Name                string
	CurrentJobsUrgent   int64
	CurrentJobsReady    int64
	CurrentJobsReserved int64
	CurrentJobsBuried   int64
	CurrentJobsDelayed  int64
	TotalJobs           int64
	CurrentUsing        int64
	CurrentWatching     int64
	CurrentWaiting      int64
	CmdDelete           int64
	CmdPauseTube        int64
	Pause               int64
	PauseTimeLeft       int64
}

func (s *TubeStats) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		key := normaliseKey(k)
		SetField(s, key, v)
	}
	return nil
}

func tubeStatsParser(res []byte) *TubeStats {
	tubeStats := make(map[string]interface{})
	fmt.Println(string(res))
	results := strings.Split(string(res), "\n")
	for _, r := range results {
		stat := strings.Split(r, ":")
		if len(stat) == 2 {
			key := stat[0]
			val := stat[1]
			if key == "name" {
				tubeStats[key] = val
			} else {
				intValue, _ := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
				fmt.Println(intValue)
				tubeStats[key] = intValue
			}
		}
	}
	result := &TubeStats{}
	result.FillStruct(tubeStats)
	return result
}
