package controllers

import (
	"k8s.io/klog/v2"
	v1 "microsoft.com/springboot-discovery-controller/api/v1"
	"net/http"
)

const (
	SUCCEEDED = "Succeeded"
	FAILED    = "Failed"
)

type Hooks[T any] struct {
	Upsert       func(target *T) error
	UpdateStatus func(target *T) error
	Delete       func(target *T) error
}

type SpringBootModel struct {
	springBootServer   *v1.SpringBootServer
	hooks              Hooks[v1.SpringBootServer]
	hooksForDiscovered Hooks[v1.SpringBootDiscovered]
	toDelete           bool
	log                klog.Logger
}

func NewModel(springBoot *v1.SpringBootServer,
	hooks Hooks[v1.SpringBootServer], hooksForDiscovered Hooks[v1.SpringBootDiscovered],
	toDelete bool,
	log klog.Logger) *SpringBootModel {
	return &SpringBootModel{
		springBootServer:   springBoot,
		hooks:              hooks,
		hooksForDiscovered: hooksForDiscovered,
		toDelete:           toDelete,
		log:                log,
	}
}

func (s *SpringBootModel) Reconcile() error {
	var err error
	if s.toDelete {
		err = s.reconcileDelete()
	} else if s.isCompleted() {
	} else {
		err = s.reconcileDiscover()
	}

	return err
}

func (s *SpringBootModel) isCompleted() bool {
	return s.springBootServer.Status.Status == SUCCEEDED
}

func (s *SpringBootModel) reconcileDelete() error {
	return s.hooksForDiscovered.Delete(s.newDiscovered())
}

func (s *SpringBootModel) reconcileDiscover() error {
	var err error
	var resp *http.Response
	discovered := s.newDiscovered()
	if resp, err = http.Get(s.springBootServer.Spec.Server); err != nil {
		return err
	} else {
		if err = s.hooksForDiscovered.Upsert(discovered); err != nil {
			return err
		} else {
			discovered.Status.StatusCode = resp.StatusCode

			if err = s.hooksForDiscovered.UpdateStatus(discovered); err != nil {
				return err
			}
		}
	}
	return err
}

func (s *SpringBootModel) markAsFailed(err error) error {
	s.springBootServer.Status.Status = FAILED
	s.springBootServer.Status.Message = err.Error()
	return s.hooks.UpdateStatus(s.springBootServer)
}

func (s *SpringBootModel) markAsReconciled() error {
	s.springBootServer.Status.Status = SUCCEEDED
	s.springBootServer.Status.Message = ""
	return s.hooks.UpdateStatus(s.springBootServer)
}

func (s *SpringBootModel) newDiscovered() *v1.SpringBootDiscovered {
	discovered := v1.SpringBootDiscovered{}
	discovered.Name = s.springBootServer.Name
	discovered.Namespace = s.springBootServer.Namespace
	discovered.Spec.Server = s.springBootServer.Spec.Server
	return &discovered

}

func (s *SpringBootModel) update() error {
	return s.hooks.Upsert(s.springBootServer)
}
