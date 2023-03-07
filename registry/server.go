package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const ServerPort = ":3000"

// 服务注册的地址，通过这个地址来注册\取消微服务
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	// Key: service name;
	// Val: Registration[]
	registrations sync.Map
}

// TODO
func (r *registry) heartbeat(freq time.Duration) {
	for {
		r.registrations.Range(func(serviceName, registration any) bool {
			go func(serviceName, registration any) {
				value, ok := registration.([]Registration)
				if !ok {
					log.Printf("type assertion failed when sending heartbeat to %s, %v is not a []Registration type, is a %v",
						serviceName, registration, reflect.TypeOf(registration))
				}
				for _, reg := range value {
					success := false
					counter := 0
				RETRY:
					res, err := http.Get(reg.HeartbeatURL)
					if err == nil || res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						success = true
					} else {
						counter++
						if counter >= 3 {
							log.Printf("Heartbeat check for %v failed after 3 retry", reg.ServiceName)
						} else {
							time.Sleep(1 * time.Second)
							goto RETRY
						}
					}
					if !success {
						r.remove(reg.ServiceURL)
					}
				}
			}(serviceName, registration)
			return true
		})
		time.Sleep(freq)
	}
}

func (r *registry) add(reg Registration) error {
	serviceName := reg.ServiceName
	registration, exist := r.registrations.Load(serviceName)
	if exist {
		value, ok := registration.([]Registration)
		if !ok {
			log.Printf("type assertion failed when adding %s, %v is not a []Registration type, is a %v",
				serviceName, registration, reflect.TypeOf(registration))
		}
		value = append(value, reg)
		registration = value
		r.registrations.Store(serviceName, registration)
	} else {
		r.registrations.Store(serviceName, []Registration{reg})
	}
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return err
}

func (r registry) notify(fullPatch patch) {
	r.registrations.Range(func(serviceName, registration any) bool {
		go func(serviceName, registration any) {
			value, ok := registration.([]Registration)
			if !ok {
				log.Printf("type assertion failed when notifying %s, %v is not a []Registration type, is a %v",
					serviceName, registration, reflect.TypeOf(registration))
			}
			for _, reg := range value {
				for _, reqService := range reg.RequiredServices {
					p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
					sendUpdate := false
					for _, added := range fullPatch.Added {
						if added.Name == reqService {
							p.Added = append(p.Added, added)
							sendUpdate = true
						}
					}
					for _, removed := range fullPatch.Removed {
						if removed.Name == reqService {
							p.Removed = append(p.Removed, removed)
							sendUpdate = true
						}
					}
					if sendUpdate {
						err := r.sendPatch(p, reg.ServiceUpdateURL)
						if err != nil {
							log.Println(err)
							return
						}
					}

				}
			}

		}(serviceName, registration)
		return true
	})
}

func (r registry) sendRequiredServices(reg Registration) error {
	var p patch
	r.registrations.Range(func(serviceName, registrations any) bool {
		value, ok := registrations.([]Registration)
		if !ok {
			log.Printf("type assertion failed when sending required services for %s, %v is not a []Registration type, is a %v",
				serviceName, registrations, reflect.TypeOf(registrations))
		}
		for _, reqServices := range reg.RequiredServices {
			if serviceName == reqServices {
				p.Added = append(p.Added, patchEntry{
					Name: value[0].ServiceName,
					URL:  value[0].ServiceURL,
				})
			}
		}
		return true
	})
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return nil
	}
	return nil

}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	var flag bool = false
	reg.registrations.Range(func(serviceName, registrations any) bool {
		value, ok := registrations.([]Registration)
		if !ok {
			log.Printf("type assertion failed when removing %s, %v is not a []Registration type, is a %v",
				url, registrations, reflect.TypeOf(registrations))
		}
		for i, regisregistration := range value {
			if regisregistration.ServiceURL == url {
				r.notify(patch{
					Removed: []patchEntry{
						{
							Name: serviceName.(ServiceName),
							URL:  regisregistration.ServiceURL,
						},
					},
				})
				value = append(value[:i], value[i+1:]...)
				if len(value) == 0 {
					reg.registrations.Delete(serviceName)
				} else {
					reg.registrations.Store(serviceName, value)
				}
				flag = true
				return false
			}
		}
		return true
	})
	if !flag {
		return fmt.Errorf("service at URL %s not found", url)
	} else {
		return nil
	}
}

// 包级变量，存放所有服务的注册信息
// 被导入时就会初始化
var reg = registry{
	registrations: sync.Map{},
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

// implement a http.Handler
type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding Service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL: %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
