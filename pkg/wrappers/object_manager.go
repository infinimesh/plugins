package wrappers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/InfiniteDevices/plugins/pkg/api"
)

type Process interface {
	Start()
	Stop()
}

type objectManagerImpl struct {
	sync.Mutex
	api                   api.Handler
	workerFactory         func(api.Object) Process
	workers               map[string]Process
	objectRefreshInterval time.Duration
	done                  chan struct{}
}

// NewObjectManager returns a new object manager
//
// It handles re-retrieval of objects in all namespaces at fixed intervals,
// routing discovered objects to new workers, while shutting down workers that
// are handling objects which are no longer part of any namespace
func NewObjectManager(handler api.Handler, workerFactory func(api.Object) Process, objectRefreshInterval time.Duration) Process {
	return &objectManagerImpl{
		api:                   handler,
		workerFactory:         workerFactory,
		workers:               map[string]Process{},
		objectRefreshInterval: objectRefreshInterval,
		done:                  make(chan struct{}),
	}
}

func (i *objectManagerImpl) Start() {
	log.Println("starting object manager...")
	i.refreshObjects() // retrieve objects upfront
	for {
		select {
		case <-i.done:
			break
		case <-time.Tick(i.objectRefreshInterval):
			if err := i.refreshObjects(); err != nil {
				log.Printf("failed to refresh objects with error: %s\n", err)
			}
		}
	}
}

func (i *objectManagerImpl) Stop() {
	i.Lock()
	defer i.Unlock()
	log.Println("stopping object manager...")
	for _, w := range i.workers {
		w.Stop()
	}
	close(i.done)
}

func (i *objectManagerImpl) refreshObjects() error {
	i.Lock()
	defer i.Unlock()
	log.Println("refreshing namespaces...")
	namespacesRes, err := i.api.GetNamespaces()
	if err != nil {
		return fmt.Errorf("failed to get namespaces: %w", err)
	}
	log.Println("refreshing objects...")
	currObjects := map[string]api.Object{}
	if namespacesRes != nil {
		for _, ns := range namespacesRes.Namespaces {
			objectsRes, err := i.api.GetObjects(ns.ID)
			if err != nil {
				return fmt.Errorf("failed to get objects: %w", err)
			}
			for _, o := range objectsRes.Objects {
				currObjects[o.UID] = o
			}
		}
	}
	i.handleNewObjects(currObjects)
	i.handleOldObjects(currObjects)
	return nil
}

func (i *objectManagerImpl) handleNewObjects(currObjects map[string]api.Object) {
	for id := range currObjects {
		if _, ok := i.workers[id]; ok {
			continue
		}
		i.workers[id] = i.workerFactory(currObjects[id])
		go i.workers[id].Start()
	}
}

func (i *objectManagerImpl) handleOldObjects(currObjects map[string]api.Object) {
	for id := range i.workers {
		if _, ok := currObjects[id]; ok {
			continue
		}
		i.workers[id].Stop()
		delete(i.workers, id)
	}
}
