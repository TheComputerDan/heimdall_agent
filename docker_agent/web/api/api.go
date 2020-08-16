package api

import (
	"encoding/json"
	_ "fmt"
	"github.com/TheComputerDan/heimdall/docker_agent/docker/connect"
	"github.com/TheComputerDan/heimdall/docker_agent/host"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
)

// GetHost returns info in the form of `Info` struct
// defined in the `host.go` file in this project.
func GetHost(w http.ResponseWriter, r *http.Request) {
	hosts := []host.Info{}
	rb := host.BuildInfo()
	hosts = append(hosts, rb)
	json.NewEncoder(w).Encode(hosts)
}

// GetContainers returns the a list of containers in the form
// of `[] types.Containers` from the docker types package.
func GetContainers(w http.ResponseWriter, r *http.Request) {
	hostContainers := connect.Containers()
	json.NewEncoder(w).Encode(hostContainers)
}

// GetImages returns a list of images on the host machine.
func GetImages(w http.ResponseWriter, r *http.Request) {
	hostImages := connect.Images()
	json.NewEncoder(w).Encode(hostImages)
}

func loadConfig() string {
	viper.SetConfigName("agent.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("docker_agent/config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	portNum := viper.GetString("server_api_port")
	return portNum
}

// Start instantiates the API and sets up the endpoints for
// consumption by the server.
func Start() {
	portNum := ":" + loadConfig()

	router := mux.NewRouter()
	router.HandleFunc("/containers", GetContainers)
	router.HandleFunc("/images", GetImages)
	router.HandleFunc("/host", GetHost)
	handlers.AllowedOrigins([]string{"*"})
	// http.ListenAndServe(":8080", router)
	http.ListenAndServe(portNum, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	// handle CORS for local testing purposes
}
