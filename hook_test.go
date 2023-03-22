package rotateloghook

import (
	"bytes"
	"github.com/sirupsen/logrus"

	"io/ioutil"
	"testing"
)

const expectedMsg = "This is the expected test message."
const unexpectedMsg = "This message should not be logged."

func TestRotateLogHook(t *testing.T) {
	log := logrus.New()
	hook, err := NewRotateLogHook("./log", "rotatelog", nil)
	if err != nil {
		t.Errorf("Unable to generate logfile due to err: %s", err)
	}
	log.Hooks.Add(hook)
	log.Info(expectedMsg)
	log.Warn(unexpectedMsg)

	fileName := "./log/rotatelog.log"
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Error while reading from log file: %s", err)
	}
	if !bytes.Contains(contents, []byte("msg=\""+expectedMsg+"\"")) {
		t.Errorf("Message read (%s) doesnt match message written (%s) for file: %s", contents, expectedMsg, fileName)
	}
	if !bytes.Contains(contents, []byte("msg=\""+unexpectedMsg+"\"")) {
		t.Errorf("Message read (%s) doesnt match message written (%s) for file: %s", contents, unexpectedMsg, fileName)
	}

}
