package log

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sendToLogger(logger *Logger, message string) {
	record := &LogRecord{
		timestamp: time.Now(),
		message:   message,
	}
	logger.push(record)
}

func benchmarkWithLoggerSizes(b *testing.B, masterMaxSize, loggerMaxSize int) {
	var (
		mockedOuput bytes.Buffer
		record      *LogRecord
	)

	master := NewMasterLogger(os.Stdout, masterMaxSize)
	master.stream = &mockedOuput
	logger := NewLogger(master, loggerMaxSize)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		record = &LogRecord{
			timestamp: time.Now(),
			message:   string("No, lamma, no!"),
		}
		logger.push(record)
	}
}

func benchmarkWritingBack(b *testing.B, masterMaxSize, loggerMaxSize int) {
	var mockedOuput bytes.Buffer

	record := &LogRecord{
		timestamp: time.Now(),
		message:   string("No, lamma, no!"),
	}

	master := NewMasterLogger(os.Stdout, masterMaxSize)
	master.stream = &mockedOuput
	logger := NewLogger(master, loggerMaxSize)

	masterStash := make([]*LogRecord, 0, masterMaxSize)
	for i := 0; i < masterMaxSize-loggerMaxSize; i++ {
		record.timestamp = randomizeTime(record.timestamp)
		masterStash = append(masterStash, record)
		record.timestamp = time.Now()
	}

	loggerStash := make([]*LogRecord, 0, loggerMaxSize)
	for i := 0; i < loggerMaxSize; i++ {
		record.timestamp = randomizeTime(record.timestamp)
		loggerStash = append(loggerStash, record)
		record.timestamp = time.Now()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		master.stash = masterStash[:]
		logger.stash = loggerStash[:]

		logger.writeBack()
	}
}

func randomizeTime(timestamp time.Time) time.Time {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return timestamp.Add(time.Duration(rand.Intn(10000)) * time.Second)
}

func TestSimpleLoggerWithMaster(t *testing.T) {
	master := NewMasterLogger(os.Stdout, 0)
	logger := NewLogger(master, 8)

	for i := 0; i < 5; i++ {
		sendToLogger(logger, fmt.Sprintf("%s", i))
		assert.Equal(t, len(logger.stash), i+1)
		assert.Equal(t, len(master.stash), 0)
	}
	logger.Close()
}

func TestOverflowedLogger(t *testing.T) {
	master := NewMasterLogger(os.Stdout, 15)
	logger := NewLogger(master, 8)

	for i := 0; i < 10; i++ {
		logger.Log(fmt.Sprintf("%s", i))
	}

	time.Sleep(3) // I know, right

	assert.True(t, len(logger.stash) < 3)
	assert.Equal(t, len(master.stash), 8)
}

func TestOverflowedMasterLogger(t *testing.T) {
	var mockedOuput bytes.Buffer
	master := NewMasterLogger(os.Stdout, 16)
	master.stream = &mockedOuput
	logger := NewLogger(master, 8)

	for i := 0; i < 25; i++ {
		sendToLogger(logger, fmt.Sprintf("%s", i))
	}
	assert.Equal(t, len(logger.stash), 1)
	assert.Equal(t, len(master.stash), 8)

	for i := 0; i < 16; i++ {
		line, err := mockedOuput.ReadBytes('\n')
		assert.True(t, len(line) > 0)
		assert.Nil(t, err)
	}
	line, err := mockedOuput.ReadBytes('\n')
	assert.False(t, len(line) > 0)
	assert.NotNil(t, err)
}

func BenchmarkPushingToM1024L256(b *testing.B) {
	benchmarkWithLoggerSizes(b, 1024, 256)
}

func BenchmarkPushingToM2048L512(b *testing.B) {
	benchmarkWithLoggerSizes(b, 2048, 512)
}

func BenchmarkPushingToM4096L1024(b *testing.B) {
	benchmarkWithLoggerSizes(b, 4096, 1024)
}

func BenchmarkPushingToM8192L2048(b *testing.B) {
	benchmarkWithLoggerSizes(b, 8192, 2048)
}

func BenchmarkPushStubbedMessage(b *testing.B) {
	benchmarkWithLoggerSizes(b, 5000000, 5000001)
}

func BenchmarkWriteBackM1024L256(b *testing.B) {
	benchmarkWritingBack(b, 1024, 256)
}

func BenchmarkWriteBackM2048L512(b *testing.B) {
	benchmarkWritingBack(b, 2048, 512)
}

func BenchmarkWriteBackM4096L1024(b *testing.B) {
	benchmarkWritingBack(b, 4096, 1024)
}

func BenchmarkWriteBackM8192L2048(b *testing.B) {
	benchmarkWritingBack(b, 8192, 2048)
}
