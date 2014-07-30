package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

// LogRecord represents single log record.
type LogRecord struct {
	timestamp time.Time
	message   string
	next      *LogRecord
}

// The logger object provides isolated and basic logger functionality just for
// the sake of making as less as possible system calls.
//
// Simple logger for each of our clients. Logs are written there message by
// message and when the size of the stash grows to the maxSize the stash is
// getting sorted and dumped into the master logger and right after the player
// logs out.
type Logger struct {
	master  *MasterLogger
	channel chan *LogRecord
	size    int
	maxSize int
	stash   []*LogRecord
}

// The master logger who is the only one could writes to the system log. All of
// his children shall point to him and dump their stashes. Writing to the
// system logs is happening only when the size of the stash is at least equals
// to maxSize or when the server intends to stop and NOT before the stash is
// sorted. Which makes it more than obvious that the maxSize of the master
// logger have to be several times bigger than the normal loggers'.
//
// Logger implements sort.Interface and the sorting is being handled by the
// standard library.
type MasterLogger struct {
	stream  io.Writer
	mutex   *sync.Mutex
	size    int
	maxSize int
	stash   []*LogRecord
}

// Constructor of the master logger. Again, there shall be only one of these.
func NewMasterLogger(maxSize int) *MasterLogger {
	logger := &MasterLogger{
		stream:  os.Stdout,
		mutex:   new(sync.Mutex),
		stash:   make([]*LogRecord, 0, maxSize),
		maxSize: maxSize,
	}
	return logger
}

// Construction of a simple logger for all of the connected clients.
func NewLogger(master *MasterLogger, maxSize int) *Logger {
	logger := &Logger{
		master:  master,
		channel: make(chan *LogRecord),
		stash:   make([]*LogRecord, 0, maxSize),
		maxSize: maxSize,
	}

	go func() {
		for record := range logger.channel {
			logger.push(record)
		}
	}()
	return logger
}

func (m *MasterLogger) Len() int {
	return m.size
}

func (m *MasterLogger) Less(i, j int) bool {
	return m.stash[i].timestamp.Before(m.stash[j].timestamp)
}

func (m *MasterLogger) Swap(i, j int) {
	m.stash[i], m.stash[j] = m.stash[j], m.stash[i]
}

// Just closes the logger inner channel. This should be done only when the user
// logs out, because once closed that channel can never be opened again... like
// the gates of hell, but you know, the other way around.
func (l *Logger) Close() {
	close(l.channel)
}

// This is the *ONLY* method that has to be used for logging from outside of
// this library. Gets the current time and runs the pushing it and the message
// through a channel (for the sake of synchronization) in another goroutine
// because we don't want each log call to block.
func (l *Logger) Log(message string) {
	go func(timestamp time.Time) {
		l.channel <- &LogRecord{
			timestamp: timestamp,
			message:   message,
		}
	}(time.Now())
}

// The actual method that pushes the records. It's very important this method
// to NOT being used outside this type, because it's completely unsafe in a
// concurrent environment.
func (l *Logger) push(record *LogRecord) {
	l.stash = append(l.stash, record)
	l.size++

	if l.size == l.maxSize {
		l.writeBack()
		l.stash = make([]*LogRecord, 0, l.maxSize)
		l.size = 0
	}
}

// Writes back all the records in the master logger and empties its stash.
func (l *Logger) writeBack() {
	l.master.mutex.Lock()
	defer l.master.mutex.Unlock()

	l.master.stash = append(l.master.stash, l.stash...)
	l.master.size += l.size
	if l.master.size >= l.master.maxSize {
		l.master.writeBack()
	}
}

// Sorts and dumps all of his stash to the default output and remains empty.
func (m *MasterLogger) writeBack() {
	sort.Sort(m)
	logger := log.New(m.stream, "", 0)

	for _, m := range m.stash {
		logger.Println(
			fmt.Sprintf(
				"[%s]: %s",
				m.timestamp.Format("2014/07/29 23:45:59"),
				m.message,
			),
		)
	}
	m.stash = make([]*LogRecord, 0, m.maxSize)
	m.size = 0
}
