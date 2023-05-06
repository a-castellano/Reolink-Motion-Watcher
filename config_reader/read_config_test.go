package config

import (
	"os"
	"testing"
)

func TestProcessConfigNoRabbitMqtt(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_rabbitmqtt/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without rabbitmqtt should fail.")
	} else {
		if err.Error() != "Fatal error config: no rabbitmq field was found." {
			t.Errorf("Error should be \"Fatal error config: no rabbitmq field was found.\" but error was '%s'.", err.Error())
		}
	}
}
