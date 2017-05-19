package logstash

import "testing"

func TestMain(t *testing.T) {
	SetMaxSize(10)
	SetOutput(".", "ra")
	WithFields(map[string]interface{}{
		"name": "tosone",
	}).Error("test")
}
