package logging

import "testing"

func TestLog(t *testing.T) {
	testLogs := []string{"test", "error on creation of music file", "something went wrong when inserting the data in DB"}
	for _, v := range testLogs {
		err := LogData(v)
		t.Error(err)
	}
}
