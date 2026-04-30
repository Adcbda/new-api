package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func GetCallbackAddress() string {
	if operation_setting.CustomCallbackAddress == "" {
		return common.BuildPublicURL(system_setting.ServerAddress, "")
	}
	return common.BuildPublicURL(operation_setting.CustomCallbackAddress, "")
}
