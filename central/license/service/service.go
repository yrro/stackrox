package service

import (
	"github.com/stackrox/rox/central/license/manager"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/grpc"
)

var (
	_ v1.LicenseServiceServer = (*service)(nil)
)

// New creates a new license service
func New(lockdownMode bool, licenseMgr manager.LicenseManager) grpc.APIService {
	return newService(lockdownMode, licenseMgr)
}
