package kics

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

func GetKicsService(config loader.Kics) interfaces.ScanServiceImpl {
	return newKicsService(config)
}
