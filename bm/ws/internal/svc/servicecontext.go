// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"ws/internal/config"
	"ws/internal/wsconn"
)

type ServiceContext struct {
	Config  config.Config
	Manager *wsconn.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Manager: wsconn.NewManager(),
	}
}
