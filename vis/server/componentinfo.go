package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type compoundTask struct {
	id, parentID         string
	kind, what, where    string
	startTime, endTime   float64
	pStartTime, pEndTime sql.NullFloat64
}

type TimeValue struct {
	Time  float64 `json:"time"`
	Value float64 `json:"value"`
}

type ComponentInfo struct {
	Name      string      `json:"name"`
	InfoType  string      `json:"info_type"`
	StartTime float64     `json:"start_time"`
	EndTime   float64     `json:"end_time"`
	Data      []TimeValue `json:"data"`
}

func httpComponentNames(w http.ResponseWriter, r *http.Request) {
	componentNames := getAllComponentNames()

	rsp, err := json.Marshal(componentNames)
	dieOnErr(err)

	_, err = w.Write([]byte(rsp))
	dieOnErr(err)
}

func httpComponentInfo(w http.ResponseWriter, r *http.Request) {

	compName := r.FormValue("where")
	infoType := r.FormValue("info_type")

	startTime, err := strconv.ParseFloat(r.FormValue("start_time"), 64)
	dieOnErr(err)

	endTime, err := strconv.ParseFloat(r.FormValue("end_time"), 64)
	dieOnErr(err)

	numDots, err := strconv.ParseInt(r.FormValue("num_dots"), 10, 32)
	dieOnErr(err)

	var compInfo *ComponentInfo
	switch infoType {
	case "ReqInCount":
		compInfo = calculateReqIn(
			compName, startTime, endTime, int(numDots))
	case "ReqCompleteCount":
		compInfo = calculateReqComplete(
			compName, startTime, endTime, int(numDots))
	case "AvgLatency":
		compInfo = calculateAvgLatency(
			compName, startTime, endTime, int(numDots))
	case "ConcurrentTask":
		compInfo = calculateTimeWeightedTaskCount(
			compName, infoType,
			startTime, endTime, int(numDots),
			func(t compoundTask) bool { return true },
			func(t compoundTask) float64 { return t.startTime },
			func(t compoundTask) float64 { return t.endTime },
		)
	case "BufferPressure":
		compInfo = calculateTimeWeightedTaskCount(
			compName, infoType,
			startTime, endTime, int(numDots),
			func(t compoundTask) bool {
				if t.kind != "req_in" {
					return false
				}

				if !t.pStartTime.Valid || !t.pEndTime.Valid {
					return false
				}

				return true
			},
			func(t compoundTask) float64 { return t.pStartTime.Float64 },
			func(t compoundTask) float64 { return t.startTime },
		)
	case "PendingReqOut":
		compInfo = calculateTimeWeightedTaskCount(
			compName, infoType,
			startTime, endTime, int(numDots),
			func(t compoundTask) bool { return t.kind == "req_out" },
			func(t compoundTask) float64 { return t.startTime },
			func(t compoundTask) float64 { return t.endTime },
		)
	default:
		log.Panicf("unknown info_type %s\n", infoType)
	}

	rsp, err := json.Marshal(compInfo)
	dieOnErr(err)

	_, err = w.Write([]byte(rsp))
	dieOnErr(err)
}

func calculateReqIn(
	compName string,
	startTime, endTime float64,
	numDots int,
) *ComponentInfo {
	info := &ComponentInfo{
		Name:      compName,
		InfoType:  "req_in",
		StartTime: startTime,
		EndTime:   endTime,
	}

	reqs := getComponentReqs(compName, startTime, endTime)

	totalDuration := endTime - startTime
	binDuration := totalDuration / float64(numDots)
	for i := 0; i < numDots; i++ {
		binStartTime := float64(i)*binDuration + startTime
		binEndTime := float64(i+1)*binDuration + startTime

		reqCount := 0
		for _, r := range reqs {
			if r.startTime > binStartTime && r.startTime < binEndTime {
				reqCount++
			}
		}

		tv := TimeValue{
			Time:  binStartTime + 0.5*binDuration,
			Value: float64(reqCount) / binDuration,
		}

		info.Data = append(info.Data, tv)
	}

	return info
}

func calculateReqComplete(
	compName string,
	startTime, endTime float64,
	numDots int,
) *ComponentInfo {
	info := &ComponentInfo{
		Name:      compName,
		InfoType:  "req_complete",
		StartTime: startTime,
		EndTime:   endTime,
	}

	reqs := getComponentReqs(compName, startTime, endTime)

	totalDuration := endTime - startTime
	binDuration := totalDuration / float64(numDots)
	for i := 0; i < numDots; i++ {
		binStartTime := float64(i)*binDuration + startTime
		binEndTime := float64(i+1)*binDuration + startTime

		reqCount := 0
		for _, r := range reqs {
			if r.endTime > binStartTime && r.endTime < binEndTime {
				reqCount++
			}
		}

		tv := TimeValue{
			Time:  binStartTime + 0.5*binDuration,
			Value: float64(reqCount) / binDuration,
		}

		info.Data = append(info.Data, tv)
	}

	return info
}

func calculateAvgLatency(
	compName string,
	startTime, endTime float64,
	numDots int,
) *ComponentInfo {
	info := &ComponentInfo{
		Name:      compName,
		InfoType:  "avg_latency",
		StartTime: startTime,
		EndTime:   endTime,
	}

	reqs := getComponentReqs(compName, startTime, endTime)

	totalDuration := endTime - startTime
	binDuration := totalDuration / float64(numDots)
	for i := 0; i < numDots; i++ {
		binStartTime := float64(i)*binDuration + startTime
		binEndTime := float64(i+1)*binDuration + startTime

		sum := 0.0
		reqCount := 0
		for _, r := range reqs {
			if r.endTime > binStartTime && r.endTime < binEndTime {
				sum += r.endTime - r.startTime
				reqCount++
			}
		}

		value := 0.0
		if reqCount > 0 {
			value = sum / float64(reqCount)
		}

		tv := TimeValue{
			Time:  binStartTime + 0.5*binDuration,
			Value: value,
		}

		info.Data = append(info.Data, tv)
	}

	return info
}

type timestamp struct {
	time    float64
	isStart bool
}

type timestamps []timestamp

func (ts timestamps) Len() int {
	return len(ts)
}

func (ts timestamps) Less(i, j int) bool {
	return ts[i].time < ts[j].time
}

func (ts timestamps) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

type taskFilter func(t compoundTask) bool
type taskTime func(t compoundTask) float64

func calculateTimeWeightedTaskCount(
	compName, infoType string,
	startTime, endTime float64,
	numDots int,
	filter taskFilter,
	increaseTime, decreaseTime taskTime,
) *ComponentInfo {
	info := &ComponentInfo{
		Name:      compName,
		InfoType:  infoType,
		StartTime: startTime,
		EndTime:   endTime,
	}

	tasks := getComponentTasks(compName, startTime, endTime)
	tasks = filterTask(tasks, filter)

	totalDuration := endTime - startTime
	binDuration := totalDuration / float64(numDots)
	for i := 0; i < numDots; i++ {
		binStartTime := float64(i)*binDuration + startTime
		binEndTime := float64(i+1)*binDuration + startTime

		tasksInBin := getTasksInBin(
			tasks,
			binStartTime, binEndTime,
			increaseTime, decreaseTime,
		)
		timestamps := taskToTimeStamps(tasksInBin, increaseTime, decreaseTime)
		avgCount := calculateAvgTaskCount(
			timestamps, binStartTime, binEndTime)

		tv := TimeValue{
			Time:  binStartTime + 0.5*binDuration,
			Value: avgCount,
		}

		info.Data = append(info.Data, tv)
	}

	return info
}

func filterTask(tasks []compoundTask, filter taskFilter) []compoundTask {
	filteredTasks := []compoundTask{}

	for _, t := range tasks {
		if filter(t) {
			filteredTasks = append(filteredTasks, t)
		}
	}

	return filteredTasks
}

// func calculateConcurrentTask(
// 	compName string,
// 	startTime, endTime float64,
// 	numDots int,
// ) *ComponentInfo {
// 	info := &ComponentInfo{
// 		Name:      compName,
// 		InfoType:  "concurrent_task",
// 		StartTime: startTime,
// 		EndTime:   endTime,
// 	}

// 	tasks := getComponentTasks(compName, startTime, endTime)

// 	totalDuration := endTime - startTime
// 	binDuration := totalDuration / float64(numDots)
// 	for i := 0; i < numDots; i++ {
// 		binStartTime := float64(i)*binDuration + startTime
// 		binEndTime := float64(i+1)*binDuration + startTime

// 		tasksInBin := getTasksInBin(tasks, binStartTime, binEndTime)
// 		timestamps := taskToTimeStamps(tasksInBin)
// 		avgCount := calculateAvgTaskCount(
// 			timestamps, binStartTime, binEndTime)

// 		tv := TimeValue{
// 			Time:  binStartTime + 0.5*binDuration,
// 			Value: avgCount,
// 		}

// 		info.Data = append(info.Data, tv)
// 	}

// 	return info
// }

// func calculateBufferPressure(
// 	compName string,
// 	startTime, endTime float64,
// 	numDots int,
// ) *ComponentInfo {
// 	info := &ComponentInfo{
// 		Name:      compName,
// 		InfoType:  "buffer_pressure",
// 		StartTime: startTime,
// 		EndTime:   endTime,
// 	}

// 	reqs := getComponentReqs(compName, startTime, endTime)

// 	totalDuration := endTime - startTime
// 	binDuration := totalDuration / float64(numDots)
// 	for i := 0; i < numDots; i++ {
// 		binStartTime := float64(i)*binDuration + startTime
// 		binEndTime := float64(i+1)*binDuration + startTime

// 		reqsInBin := getReqBufferredInBin(reqs, binStartTime, binEndTime)
// 		timestamps := reqBufferredToTimeStamps(reqsInBin)
// 		avgCount := calculateAvgTaskCount(
// 			timestamps, binStartTime, binEndTime)

// 		tv := TimeValue{
// 			Time:  binStartTime + 0.5*binDuration,
// 			Value: avgCount,
// 		}

// 		info.Data = append(info.Data, tv)
// 	}

// 	return info
// }

// getReqBufferredInBin filters the requests. It returns the requests that has /// been a buffer between the binStart and binEnd time.
// func getReqBufferredInBin(
// 	reqs []compoundTask,
// 	binStart, binEnd float64,
// ) []compoundTask {
// 	reqsInBin := []compoundTask{}

// 	for _, t := range reqs {
// 		if isReqBufferredOverlapsWithBin(t, binStart, binEnd) {
// 			reqsInBin = append(reqsInBin, t)
// 		}
// 	}

// 	return reqsInBin
// }

func calculateAvgTaskCount(
	timestamps timestamps,
	binStartTime, binEndTime float64,
) float64 {
	var count int
	var timeByCount float64
	prevTime := binStartTime

	for _, ts := range timestamps {
		if ts.time < binStartTime {
			if ts.isStart {
				count++
			} else {
				count--
			}
			continue
		} else if ts.time >= binEndTime {
			break
		} else {
			duration := ts.time - prevTime
			if duration < 0 {
				panic("duration is smaller than 0")
			}
			timeByCount += duration * float64(count)
			prevTime = ts.time

			if ts.isStart {
				count++
			} else {
				count--
			}
		}
	}

	duration := binEndTime - prevTime
	timeByCount += duration * float64(count)

	avgCount := timeByCount / (binEndTime - binStartTime)

	return avgCount
}

func taskToTimeStamps(
	tasks []compoundTask,
	taskStart, taskEnd taskTime,
) []timestamp {
	var timestamps timestamps

	for _, t := range tasks {
		timestampStart := timestamp{
			time:    taskStart(t),
			isStart: true,
		}

		timestampEnd := timestamp{
			time: taskEnd(t),
		}

		timestamps = append(timestamps, timestampStart, timestampEnd)
	}

	sort.Sort(timestamps)

	return timestamps
}

// func reqBufferredToTimeStamps(tasks []compoundTask) []concurrentTaskTimestamp {
// 	var timestamps timestamps

// 	for _, t := range tasks {
// 		timestampStart := concurrentTaskTimestamp{
// 			time:    t.pStartTime.Float64,
// 			isStart: true,
// 		}

// 		timestampEnd := concurrentTaskTimestamp{
// 			time: t.startTime,
// 		}

// 		timestamps = append(timestamps, timestampStart, timestampEnd)
// 	}

// 	sort.Sort(timestamps)

// 	return timestamps
// }

func getTasksInBin(
	tasks []compoundTask,
	binStart, binEnd float64,
	taskStart, taskEnd taskTime,
) (tasksInBin []compoundTask) {
	for _, t := range tasks {
		if isTaskOverlapsWithBin(t, binStart, binEnd, taskStart, taskEnd) {
			tasksInBin = append(tasksInBin, t)
		}
	}

	return tasksInBin
}

func isTaskOverlapsWithBin(
	t compoundTask,
	binStart, binEnd float64,
	taskStart, taskEnd taskTime,
) bool {
	if taskEnd(t) < binStart {
		return false
	}

	if taskStart(t) > binEnd {
		return false
	}

	return true
}

// func isReqBufferredOverlapsWithBin(
// 	r compoundTask,
// 	binStart, binEnd float64,
// ) bool {
// 	if !r.pStartTime.Valid {
// 		return false
// 	}

// 	if r.pStartTime.Float64 > r.startTime {
// 		panic("never")
// 	}

// 	if r.pStartTime.Float64 < binStart {
// 		return false
// 	}

// 	if r.startTime > binEnd {
// 		return false
// 	}

// 	return true
// }

var prepareSelectCompReqSqlStmtOnce sync.Once
var selectCompReqSqlStmt *sql.Stmt

func getComponentReqs(name string, startTime, endTime float64) []compoundTask {
	stmt := getSelectComponentReqStmt()
	rows, err := stmt.Query(name, startTime, endTime)
	dieOnErr(err)

	tasks := []compoundTask{}
	for rows.Next() {
		// var start, end float64
		t := compoundTask{}

		err := rows.Scan(
			&t.id, &t.parentID,
			&t.kind, &t.what, &t.where,
			&t.startTime, &t.endTime,
			&t.pStartTime, &t.pEndTime,
		)
		dieOnErr(err)

		tasks = append(tasks, t)
	}

	return tasks
}

func getSelectComponentReqStmt() *sql.Stmt {
	prepareSelectCompReqSqlStmtOnce.Do(func() {
		sqlStr := `
			SELECT
				t.task_id,
				t.parent_id,
				t.kind,
				t.what,
				t.location,
				t.start_time,
				t.end_time,
				pt.start_time,
				pt.end_time
			FROM trace t
			LEFT JOIN trace pt
			ON t.parent_id = pt.task_id
			WHERE t.kind='req_in' 
				AND t.location=?
				AND t.end_time>?
				AND t.start_time<?
		`

		stmt, err := db.Prepare(sqlStr)
		dieOnErr(err)

		selectCompReqSqlStmt = stmt
	})

	return selectCompReqSqlStmt
}

var prepareSelectCompTaskSqlStmtOnce sync.Once
var selectCompTaskSqlStmt *sql.Stmt

func getComponentTasks(name string, startTime, endTime float64) []compoundTask {
	stmt := getSelectComponentTaskStmt()
	rows, err := stmt.Query(name, startTime, endTime)
	dieOnErr(err)

	tasks := []compoundTask{}
	for rows.Next() {
		// var start, end float64
		t := compoundTask{}

		err := rows.Scan(
			&t.id, &t.parentID,
			&t.kind, &t.what, &t.where,
			&t.startTime, &t.endTime,
			&t.pStartTime, &t.pEndTime,
		)
		dieOnErr(err)

		tasks = append(tasks, t)
	}

	return tasks
}

func getSelectComponentTaskStmt() *sql.Stmt {
	prepareSelectCompTaskSqlStmtOnce.Do(func() {
		sqlStr := `
		SELECT
			t.task_id,
			t.parent_id,
			t.kind,
			t.what,
			t.location,
			t.start_time,
			t.end_time,
			pt.start_time,
			pt.end_time
		FROM trace t
		LEFT JOIN trace pt
		ON t.parent_id = pt.task_id
		WHERE t.location=?
			AND t.end_time>?
			AND t.start_time<?
		`

		stmt, err := db.Prepare(sqlStr)
		dieOnErr(err)

		selectCompTaskSqlStmt = stmt
	})

	return selectCompTaskSqlStmt
}

func getAllComponentNames() []string {
	rows, err := db.Query("SELECT DISTINCT location FROM trace")
	dieOnErr(err)

	names := []string{}

	for rows.Next() {
		componentName := ""
		err := rows.Scan(&componentName)
		dieOnErr(err)

		if componentName == "" {
			continue
		}

		names = append(names, componentName)
	}

	return names
}
