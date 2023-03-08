# Distributed :globe_with_meridians:
Distributed is a Go-based microservices architecture system that provides service registration, discovery, and dependency update. You can put your custome micro services under the cmd directory to run the system.
## :gear: Prerequisites
Go v1.18 or higher

Docker (optional)
## :rocket: Getting Started
Clone this repository:
``` bash
git clone https://github.com/your_username/distributed.git
```
Navigate to the cloned directory:
``` bash
cd distributed
```
Start the micro services by navigating to their respective directories under /cmd and run:
``` bash
go run main.go
```
For example, to start the registry service:
``` bash
cd cmd/registryservice
go run main.go
```
Use the API endpoints provided by the `package service` to interact with the system.

For example, to register a service with the `registry`:
``` golang
host, port := "localhost", "4000"
serviceAddr := fmt.Sprintf("http://%s:%s", host, port)
r := registry.Registration{
	ServiceName:      registry.LogService,
	ServiceURL:       serviceAddr,
	RequiredServices: make([]registry.ServiceName, 0),
	ServiceUpdateURL: serviceAddr + "/services",
	HeartbeatURL:     serviceAddr + "/heartbeat",
}
ctx, err := service.Start(
	context.Background(),
	host,
	port,
	r,
	log.RegisterHandlers,
)
```
## :whale: Running with Docker
Distributed can also be run using Docker.

`TBD`
## :hammer: Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## :page_with_curl: License
This project is licensed under the MIT License - see the `LICENSE` file for details.