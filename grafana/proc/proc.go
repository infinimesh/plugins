package proc

import (
	"fmt"
	"log"
	"time"

	"github.com/InfiniteDevices/plugins/grafana/api"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
)

type Proc interface {
	Run(time.Duration)
}

type procImpl struct {
	api         api.Handler
	redisClient *redistimeseries.Client
	device      map[string]*objectWorker
}

func New(handler api.Handler, redisClient *redistimeseries.Client) Proc {
	return &procImpl{
		api:         handler,
		redisClient: redisClient,
		device:      map[string]*objectWorker{},
	}
}

func (p *procImpl) Run(tick time.Duration) {
	log.Println("starting processor...")
	p.refreshDevices()
	for range time.Tick(tick) {
		if err := p.refreshDevices(); err != nil {
			log.Printf("failed to refresh devices with error: %s\n", err)
		}
	}
}

func (p *procImpl) refreshDevices() error {
	log.Println("refreshing devices...")
	namespacesRes, err := p.api.GetNamespaces()
	if err != nil {
		return fmt.Errorf("failed to get namespaces: %w", err)
	}

	currDevices := map[string]api.Object{}
	if namespacesRes != nil {
		for _, ns := range namespacesRes.Namespaces {
			objectsRes, err := p.api.GetObjects(ns.ID)
			if err != nil {
				return fmt.Errorf("failed to get objects: %w", err)
			}
			for _, o := range objectsRes.Objects {
				currDevices[o.UID] = o
			}
		}
	}

	// handle new devices
	for id := range currDevices {
		if _, ok := p.device[id]; ok {
			continue
		}
		p.device[id] = &objectWorker{
			obj:         currDevices[id],
			api:         p.api,
			redisClient: p.redisClient,
		}
		go p.device[id].Start()
	}

	// handle old devices
	for id := range p.device {
		if _, ok := currDevices[id]; ok {
			continue
		}
		p.device[id].Stop()
		delete(p.device, id)
	}

	return nil
}
