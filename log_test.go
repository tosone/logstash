package logstash

import "testing"

func TestMain(t *testing.T) {
	SetMaxSize(10)
	SetOutput(".", "ra")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
	WithFields(Fields{"name": "tosone", "address": "here"}).Error("test")
}
