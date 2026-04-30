package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/stretchr/testify/require"
)

func TestGetCallbackAddressIncludesAppBasePath(t *testing.T) {
	originalBasePath := common.AppBasePath
	originalServerAddress := system_setting.ServerAddress
	originalCustomCallbackAddress := operation_setting.CustomCallbackAddress
	common.AppBasePath = "/new-api"
	system_setting.ServerAddress = "https://example.com"
	operation_setting.CustomCallbackAddress = ""
	t.Cleanup(func() {
		common.AppBasePath = originalBasePath
		system_setting.ServerAddress = originalServerAddress
		operation_setting.CustomCallbackAddress = originalCustomCallbackAddress
	})

	require.Equal(t, "https://example.com/new-api", GetCallbackAddress())

	operation_setting.CustomCallbackAddress = "https://callback.example.com/new-api"
	require.Equal(t, "https://callback.example.com/new-api", GetCallbackAddress())
}
