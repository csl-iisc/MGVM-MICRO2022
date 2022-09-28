package tracing

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/akita/akita"

	// Need to use MySQL connections.
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
	"github.com/tebeka/atexit"
)

// MySQLTracer is a task tracer that can store the tasks into a MySQL database.
type MySQLTracer struct {
	username         string
	password         string
	ipAddress        string
	port             int
	dbName           string
	db               *sql.DB
	tracingTasks     map[string]Task
	tasksToWriteToDB []Task
	batchSize        int
	startTime        akita.VTimeInSec
	endTime          akita.VTimeInSec
}

// Init establishes a connection to MySQL and creates a database.
func (t *MySQLTracer) Init() {
	t.getCredentials()
	t.connect()
	t.createDatabase()
}

func (t *MySQLTracer) getCredentials() {
	t.username = os.Getenv("AKITA_TRACE_USERNAME")
	if t.username == "" {
		panic(`trace username is not set, use environment variable AKITA_TRACE_USERNAME to set it.`)
	}

	t.password = os.Getenv("AKITA_TRACE_PASSWORD")
	t.ipAddress = os.Getenv("AKITA_TRACE_IP")
	if t.ipAddress == "" {
		t.ipAddress = "127.0.0.1"
	}

	portString := os.Getenv("AKITA_TRACE_PORT")
	if portString == "" {
		portString = "3306"
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}
	t.port = port
}

func (t *MySQLTracer) connect() {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		t.username, t.password, t.ipAddress, t.port)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		panic(err)
	}

	t.db = db
}

func (t *MySQLTracer) connectWithDBName() {
	err := t.db.Close()
	if err != nil {
		panic(err)
	}

	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		t.username, t.password, t.ipAddress, t.port, t.dbName)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		panic(err)
	}

	t.db = db
	t.db.SetMaxOpenConns(10000)
	t.db.SetMaxIdleConns(100)
	t.db.SetConnMaxLifetime(5 * time.Second)
}

func (t *MySQLTracer) createDatabase() {
	dbName := "akita_trace_" + xid.New().String()
	t.dbName = dbName
	log.Printf("Trace is Collected in Database: %s\n", dbName)

	t.mustExecute("CREATE DATABASE " + dbName)
	t.mustExecute("USE " + dbName)

	t.createTable()
}

func (t *MySQLTracer) createTable() {
	t.mustExecute(`
		create table trace
		(
			_id        int auto_increment
				primary key,
			task_id    varchar(200) null,
			parent_id  varchar(200) null,
			kind       varchar(100) null,
			what       varchar(100) null,
			location   varchar(100) null,
			start_time float       null,
			end_time   float       null
		);
	`)

	t.mustExecute(`
        ALTER TABLE trace ENGINE=MyISAM;
	`)

	t.mustExecute(`
		create index trace_end_time_index
			on trace (end_time) USING BTREE;
	`)

	t.mustExecute(`
		create index trace_task_id_uindex
			on trace (task_id);
	`)

	t.mustExecute(`
		create index trace_kind_index
			on trace (kind);
	`)

	t.mustExecute(`
		create index trace_start_time_index
			on trace (start_time) USING BTREE;
	`)

	t.mustExecute(`
		create index trace_what_index
			on trace (what);
	`)

	t.mustExecute(`
		create index trace_location_index
			on trace (location);
	`)

	t.mustExecute(`
		create index trace_parent_id_index
			on trace (parent_id);
	`)
}

func (t *MySQLTracer) mustExecute(query string) sql.Result {
	res, err := t.db.Exec(query)
	if err != nil {
		panic(err)
	}
	return res
}

// StartTask marks the start of a task.
func (t *MySQLTracer) StartTask(task Task) {
	if t.endTime > 0 && task.StartTime > t.endTime {
		return
	}

	t.tracingTasks[task.ID] = task
}

// StepTask marks a milestone during the execution of a task.
func (t *MySQLTracer) StepTask(task Task) {
	// Do nothing for now
	//t.lock.Lock()
	original, ok := t.tracingTasks[task.ID]
	if !ok {
		panic("oh no!")
	}
	taskType := task.Steps[0].What
	original.What = original.What + "-" + taskType
	t.tracingTasks[task.ID] = original
	//fmt.Println(taskType, original.What)
}

// EndTask writes the task into the database.
func (t *MySQLTracer) EndTask(task Task) {
	if t.startTime > 0 && task.EndTime < t.startTime {
		delete(t.tracingTasks, task.ID)
		return
	}

	originalTask := t.tracingTasks[task.ID]
	originalTask.EndTime = task.EndTime
	originalTask.Detail = nil
	delete(t.tracingTasks, task.ID)
	// fmt.Println(originalTask, originalTask.Kind, originalTask.What)
	t.tasksToWriteToDB = append(t.tasksToWriteToDB, originalTask)
	if len(t.tasksToWriteToDB) > t.batchSize {
		t.flushToDB()
	}
}

func (t *MySQLTracer) flushToDB() {
	sqlStr := `INSERT INTO trace VALUES`
	vals := []interface{}{}

	for i := range t.tasksToWriteToDB {
		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals,
			0,
			t.tasksToWriteToDB[i].ID,
			t.tasksToWriteToDB[i].ParentID,
			t.tasksToWriteToDB[i].Kind,
			t.tasksToWriteToDB[i].What,
			t.tasksToWriteToDB[i].Where,
			t.tasksToWriteToDB[i].StartTime,
			t.tasksToWriteToDB[i].EndTime,
		)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// fmt.Println(sqlStr)
	stmt, err := t.db.Prepare(sqlStr)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(vals...)
	if err != nil {
		panic(err)
	}

	err = stmt.Close()
	if err != nil {
		panic(err)
	}

	t.tasksToWriteToDB = nil
}

// NewMySQLTracerWithTimeRange creates a MySQLTracer which can only trace the
// tasks that at least partially overlaps with the given start and end time. If
// the start time is negative, the tracer will start tracing at the beginning of
// the simulation. If the end time is negative, the tracer will not stop tracing
// until the end of the simulation.
func NewMySQLTracerWithTimeRange(
	startTime, endTime akita.VTimeInSec,
) *MySQLTracer {
	if startTime >= 0 && endTime >= 0 {
		if startTime >= endTime {
			panic("startTime cannot be greater than endTime")
		}
	}

	t := &MySQLTracer{
		startTime:    startTime,
		endTime:      endTime,
		tracingTasks: make(map[string]Task),
		batchSize:    4000,
	}

	atexit.Register(func() { t.flushToDB() })

	return t
}

// NewMySQLTracer returns a new MySQLTracer.
func NewMySQLTracer() *MySQLTracer {
	t := &MySQLTracer{
		startTime:    -1,
		endTime:      -1,
		tracingTasks: make(map[string]Task),
		batchSize:    4000,
	}

	atexit.Register(func() { t.flushToDB() })

	return t
}
