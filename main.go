package main

import (
	"github.com/govindansriram/qldriver"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

var driver qldriver.PublisherClient

type config struct {
	Name             string `yaml:"name"`
	Password         string `yaml:"password"`
	Address          string `yaml:"address"`
	Port             uint16 `yaml:"port"`
	MaxConnections   uint16 `yaml:"maxPublisherConnections"`
	MaxIoTimeSeconds uint16 `yaml:"maxIoTimeSeconds"`
}

func readConfig(fileData []byte) (qldriver.PublisherClient, error) {
	structure := config{}
	err := yaml.Unmarshal(fileData, &structure)

	if err != nil {
		return qldriver.PublisherClient{}, err
	}

	return qldriver.NewPublisherClient(
		structure.Name,
		structure.Password,
		structure.MaxConnections,
		structure.MaxIoTimeSeconds,
		structure.Port,
		structure.Address)
}

func getDriver() qldriver.PublisherClient {
	return driver
}

func main() {

	fPath := os.Args[1]
	data, err := os.ReadFile(fPath)

	if err != nil {
		log.Fatal("could not read configuration file")
	}

	qDriver, err := readConfig(data)

	if err != nil {
		log.Fatal(err)
	}

	driver = qDriver

	apiKey := os.Getenv("APIKEY")

	if apiKey == "" {
		log.Fatal("apikey was not provided")
	}

	handler := enqueueHandler{
		key: apiKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handler)

	err = http.ListenAndServe(":5000", mux)

	if err != nil {
		log.Fatal(err)
	}
}
