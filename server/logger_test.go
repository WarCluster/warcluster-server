package server

import (
	"bytes"
	"fmt"
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

func TestSimpleLoggerWithMaster(t *testing.T) {
	master := NewMasterLogger(0)
	logger := NewLogger(master, 8)

	for i := 0; i < 5; i++ {
		sendToLogger(logger, fmt.Sprintf("%s", i))
		assert.Equal(t, logger.size, i+1)
		assert.Equal(t, master.size, 0)
	}
	logger.Close()
}

func TestOverflowedLogger(t *testing.T) {
	master := NewMasterLogger(15)
	logger := NewLogger(master, 8)

	for i := 0; i < 10; i++ {
		logger.Log(fmt.Sprintf("%s", i))
	}

	time.Sleep(3) // I know, right

	assert.True(t, logger.size < 3)
	assert.Equal(t, master.size, 8)
}

func TestOverflowedMasterLogger(t *testing.T) {
	var mockedOuput bytes.Buffer
	master := NewMasterLogger(16)
	master.stream = &mockedOuput
	logger := NewLogger(master, 8)

	for i := 0; i < 25; i++ {
		sendToLogger(logger, fmt.Sprintf("%s", i))
	}
	assert.Equal(t, logger.size, 1)
	assert.Equal(t, master.size, 8)

	for i := 0; i < 16; i++ {
		line, err := mockedOuput.ReadBytes('\n')
		assert.True(t, len(line) > 0)
		assert.Nil(t, err)
	}
	line, err := mockedOuput.ReadBytes('\n')
	assert.False(t, len(line) > 0)
	assert.NotNil(t, err)
}
