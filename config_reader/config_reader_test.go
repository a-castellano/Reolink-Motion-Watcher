package config

import (
	"os"
	"strings"
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

func TestProcessConfigNoRabbitMqttPort(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_rabbitmqtt_port/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without rabbitmqtt host should fail.")
	} else {
		if err.Error() != "Fatal error config: no rabbitmq port was found." {
			t.Errorf("Error should be \"Fatal error config: no rabbitmq port was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoAlarmManager(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_alarmmanager/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without alarmmanager should fail.")
	} else {
		if err.Error() != "Fatal error config: no alarmmanager field was found." {
			t.Errorf("Error should be \"Fatal error config: no alarmmanager field was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoAlarmManagerDeviceId(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_alarmmanager_device_id/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without alarmmanager deviceid should fail.")
	} else {
		if err.Error() != "Fatal error config: no alarmManager deviceid was found." {
			t.Errorf("Error should be \"Fatal error config: no alarmManager deviceid was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebcams(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_webcams/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without webcams should fail.")
	} else {
		if err.Error() != "Fatal error config: no webcams field was found." {
			t.Errorf("Error should be \"Fatal error config: no webcams field was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWebcamNameRepeated(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_webcam_name_repeated/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with webcam name repeated should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam1 name is repeated." {
			t.Errorf("Error should be \"Fatal error config: webcam cam1 name is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWebcamIPRepeated(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_webcam_ip_repeated/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with webcam ip repeated should fail.")
	} else {
		if !strings.Contains(err.Error(), "ip is repeated.") {
			t.Errorf("Error should contain \"ip is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWebcamIDRepeated(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_webcam_id_repeated/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with webcam id repeated should fail.")
	} else {
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Error should contain \"already exists\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebcamsEntries(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_webcams_entries/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without webcams entries should fail.")
	} else {
		if err.Error() != "Fatal error config: no webcams were found." {
			t.Errorf("Error should be \"Fatal error config: no webcams were found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoRedis(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_redis/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without redis should fail.")
	} else {
		if err.Error() != "Fatal error config: no redis field was found." {
			t.Errorf("Error should be \"Fatal error config: no redis field was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigRepeatedQueueName(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_same_queue_name/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig with repeated queue name should fail.")
	} else {
		if err.Error() != "Fatal error config: rabbitmq motion_queue and video_queue cannot be the same." {
			t.Errorf("Error should be \"Fatal error config: rabbitmq motion_queue and video_queue cannot be the same.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigOK(t *testing.T) {
	os.Setenv("MOTION_WATCHER_CONFIG_FILE_LOCATION", "./config_files_test/config_ok/")
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig method with right config should not fail.")
	}
}
