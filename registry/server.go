package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"

// 服务注册的地址，通过这个地址来注册微服务
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	// to protect the following slice
	mutex        *sync.Mutex
	registration []Registration
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registration = append(r.registration, reg)
	r.mutex.Unlock()
	return nil
}

var reg = registry{
	registration: make([]Registration, 0),
	mutex:        new(sync.Mutex),
}

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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
