package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/akita/util/tracing"
)

func httpTrace(w http.ResponseWriter, r *http.Request) {
	kindLabel := r.FormValue("kind")
	idLabel := r.FormValue("id")
	parentIDLabel := r.FormValue("parentid")
	whereLabel := r.FormValue("where")
	startTimeLabel := r.FormValue("starttime")
	endTimeLabel := r.FormValue("endtime")

	sqlStr := `
		SELECT 
			task_id, 
			parent_id, 
			kind, 
			what, 
			location, 
			start_time, 
			end_time
		FROM trace 
		WHERE 
	`
	vals := []interface{}{}

	firstCondition := true

	if idLabel != "" {
		sqlStr += "task_id=?"
		vals = append(vals, idLabel)
		firstCondition = false
	}

	if kindLabel != "" {
		if !firstCondition {
			sqlStr += " AND "
		}
		sqlStr += "kind=?"
		vals = append(vals, kindLabel)
		firstCondition = false
	}

	if parentIDLabel != "" {
		if !firstCondition {
			sqlStr += " AND "
		}
		sqlStr += "parent_id=?"
		vals = append(vals, parentIDLabel)
		firstCondition = false
	}

	if whereLabel != "" {
		if !firstCondition {
			sqlStr += " AND "
		}
		sqlStr += "location=?"
		vals = append(vals, whereLabel)
		firstCondition = false
	}

	if startTimeLabel != "" {
		startTime, err := strconv.ParseFloat(startTimeLabel, 64)
		dieOnErr(err)

		if !firstCondition {
			sqlStr += " AND "
		}
		sqlStr += "end_time>?"
		vals = append(vals, startTime)
		firstCondition = false
	}

	if endTimeLabel != "" {
		endTime, err := strconv.ParseFloat(endTimeLabel, 64)
		dieOnErr(err)

		if !firstCondition {
			sqlStr += " AND "
		}
		sqlStr += "start_time<?"
		vals = append(vals, endTime)
		// firstCondition = false
	}

	stmt, err := db.Prepare(sqlStr)
	dieOnErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(vals...)
	dieOnErr(err)

	tasks := make([]tracing.Task, 0)
	for rows.Next() {
		task := tracing.Task{}
		err := rows.Scan(
			&task.ID,
			&task.ParentID,
			&task.Kind,
			&task.What,
			&task.Where,
			&task.StartTime,
			&task.EndTime,
		)
		dieOnErr(err)

		tasks = append(tasks, task)
	}

	rsp, err := json.Marshal(tasks)
	dieOnErr(err)

	_, err = w.Write([]byte(rsp))
	dieOnErr(err)
}
