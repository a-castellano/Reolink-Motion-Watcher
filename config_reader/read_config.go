package config

import (
	"errors"
	"net"
	"reflect"

	webcam "github.com/a-castellano/reolink-manager/webcam"
	viperLib "github.com/spf13/viper"
)

type Rabbitmq struct {
	Host     string
	Port     int
	User     string
	Password string
}

type AlarmManager struct {
	Host     string
	Port     int
	DeviceId string
}

type Config struct {
	Rabbitmq     Rabbitmq
	AlarmManager AlarmManager
	Webcams      map[string]webcam.Webcam
}

func contains(keys []string, keyName string) bool {
	for _, v := range keys {
		if v == keyName {
			return true
		}
	}

	return false
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "MOTION_WATCHER_CONFIG_FILE_LOCATION"

	requiredVariables := []string{"rabbitmq", "alarmmanager", "webcams"}
	rabbitmqRequiredVariables := []string{"host", "port", "user", "password"}
	webcamRequiredVariables := []string{"ip", "user", "password", "name"}
	alarmManagerRequiredVariables := []string{"host", "port", "deviceid"}

	viper := viperLib.New()

	//Look for config file location defined as env var
	viper.BindEnv(envVariable)
	configFileLocation = viper.GetString(envVariable)
	if configFileLocation == "" {
		// Get config file from default location
		return config, errors.New(errors.New("Environment variable MOTION_WATCHER_CONFIG_FILE_LOCATION is not defined.").Error())
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configFileLocation)

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.New(errors.New("Fatal error reading config file: ").Error() + err.Error())
	}

	for _, requiredVariable := range requiredVariables {
		if !viper.IsSet(requiredVariable) {
			return config, errors.New("Fatal error config: no " + requiredVariable + " field was found.")
		}
	}

	for _, rabbitmqVariable := range rabbitmqRequiredVariables {
		if !viper.IsSet("rabbitmq." + rabbitmqVariable) {
			return config, errors.New("Fatal error config: no rabbitmq " + rabbitmqVariable + " was found.")
		}
	}

	for _, alarmManagerVariable := range alarmManagerRequiredVariables {
		if !viper.IsSet("alarmmanager." + alarmManagerVariable) {
			return config, errors.New("Fatal error config: no alarmManager " + alarmManagerVariable + " was found.")
		}
	}

	webcams := make(map[string]webcam.Webcam)
	readedWebCamIDs := make(map[string]bool)
	readedWebCamNames := make(map[string]bool)
	readedWebCamIPs := make(map[string]bool)
	readedWebcams := viper.GetStringMap("webcams")
	for webcamID, webcamInfo := range readedWebcams {
		webCamName := "NoName"
		webcamInfoValue := reflect.ValueOf(webcamInfo)
		var newWebcam webcam.Webcam
		if webcamInfoValue.Kind() != reflect.Map {
			return config, errors.New("Fatal error config: webcam " + webcamID + " not a map.")
		} else {

			if _, ok := readedWebCamIDs[webcamID]; ok {
				return config, errors.New("Fatal error config: webcam " + webcamID + " is repeated.")
			} else {

				webcamInfoValueMap := webcamInfoValue.Interface().(map[string]interface{})

				keys := make([]string, 0, len(webcamInfoValueMap))
				for key_name := range webcamInfoValueMap {
					keys = append(keys, key_name)
				}
				for _, requiredWebcamKey := range webcamRequiredVariables {
					if !contains(keys, requiredWebcamKey) {
						return config, errors.New("Fatal error config: webcam " + webcamID + " has no " + requiredWebcamKey + ".")
					} else {
						if requiredWebcamKey == "ip" {
							newWebcam.IP = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
							if net.ParseIP(newWebcam.IP) == nil {
								return config, errors.New("Fatal error config: webcam " + webcamID + " ip is invalid.")
							} else {
								if _, ok := readedWebCamIPs[newWebcam.IP]; ok {
									return config, errors.New("Fatal error config: webcam " + webcamID + " ip is repeated.")
								} else {
									readedWebCamIPs[newWebcam.IP] = true
								}
							}
						} else {
							if requiredWebcamKey == "user" {
								newWebcam.User = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
							} else {
								if requiredWebcamKey == "password" {
									newWebcam.Password = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
								} else {
									if requiredWebcamKey == "name" {
										webCamName = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
										if _, ok := readedWebCamNames[webCamName]; ok {
											return config, errors.New("Fatal error config: webcam " + webCamName + " name is repeated.")
										} else {
											readedWebCamNames[webCamName] = true
										}
									}
								}
							}
						}
					}
				}

				webcams[webCamName] = newWebcam
			}
		}
	}

	if len(webcams) == 0 {
		return config, errors.New("Fatal error config: no webcams were found.")
	}

	rabbitmqConfig := Rabbitmq{Host: viper.GetString("rabbitmq.host"), Port: viper.GetInt("rabbitmq.port"), User: viper.GetString("rabbitmq.user"), Password: viper.GetString("rabbitmq.password")}

	alarmManagerConfig := AlarmManager{Host: viper.GetString("alarmmanager.host"), Port: viper.GetInt("alarmmanager.port"), DeviceId: viper.GetString("alarmmanager.deviceid")}

	config.Rabbitmq = rabbitmqConfig
	config.AlarmManager = alarmManagerConfig
	config.Webcams = webcams

	return config, nil
}
