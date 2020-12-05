package proc

import (
	"log"

	"github.com/InfiniteDevices/plugins/grafana/api"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
)

// dataKeys are the keys in shadow.reported.data that are numeric
var dataKeys = []string{
	"ac_realpower",
	"apparent_power",
	"day_yield",
	"dc_inc_volt_A",
	"dc_input",
	"grid_current_ampere",
	"grid_freq",
	"grid_phase_A",
	"grid_phase_B",
	"grid_phase_C",
}

type objectWorker struct {
	obj         api.Object
	api         api.Handler
	done        chan struct{}
	redisClient *redistimeseries.Client
}

func (w *objectWorker) Start() {
	createOpts := redistimeseries.CreateOptions{
		Labels: map[string]string{
			"uid":  w.obj.UID,
			"name": w.obj.Name,
			"kind": w.obj.Kind,
		},
	}

	for _, k := range dataKeys {
		err := w.redisClient.CreateKeyWithOptions(k+":"+w.obj.UID, createOpts)
		if err != nil {
			log.Printf("error creating redis key %s: %s\n", k, err)
		}
	}

	ch, err := w.api.GetDevicesStateStream(w.obj.UID)
	if err != nil {
		log.Printf("error on get devices state stream: %s\n", err)
	}

	for {
		select {
		case <-w.done:
			return
		case state := <-ch:
			for _, k := range dataKeys {
				if state == nil || state.Result.ReportedState.Data == nil || state.Result.ReportedState.Data[k] == nil {
					log.Printf("received nil data for object %v and key %s", w.obj, k)
					continue
				}
				v, ok := state.Result.ReportedState.Data[k].(float64)
				if !ok {
					log.Printf("invalid data type found for object %v and key %s", w.obj, k)
					continue
				}
				_, err = w.redisClient.AddAutoTsWithOptions(k+":"+w.obj.UID, v, createOpts)
				if err != nil {
					log.Printf("failed to add time series item: %s\n", err)
					continue
				}
				log.Printf("added time series item: object=%v key=%s\n", w.obj, k)
			}
		}
	}
}

func (w *objectWorker) Stop() {
	close(w.done)
}
