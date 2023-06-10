package base

import "fmt"

const (
	ConfigPathFromService = "../../data/config/"
)

func GetConfigFilePathByService(service string) string {
	return fmt.Sprintf("%s/%s/%s.yaml", ConfigPathFromService, service, service)
}
