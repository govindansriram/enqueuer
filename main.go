package main

import (
	"errors"
	"fmt"
	"github.com/govindansriram/qldriver"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"strings"
)

const localIp = "0.0.0.0"

type config struct {
	Name             string `yaml:"name"`
	Password         string `yaml:"password"`
	Apikey           string `yaml:"apikey"`
	MaxConnections   uint16 `yaml:"maxPublisherConnections"`
	MaxIoTimeSeconds uint16 `yaml:"maxIoTimeSeconds"`
	Port             uint16 `yaml:"port"`
	QliteAddress     string `yaml:"qliteAddress"`
	QlitePort        uint16 `yaml:"qlitePort"`
}

func readConfig(fileData []byte) (driver qldriver.PublisherClient, key, address string, err error) {
	structure := config{}
	err = yaml.Unmarshal(fileData, &structure)

	if err != nil {
		return
	}

	key = structure.Apikey

	if key == "" {
		err = errors.New("key is empty")
		return
	}

	cond1 := strings.Contains(strings.ToLower(structure.QliteAddress), "localhost")
	cond2 := strings.Contains(structure.QliteAddress, localIp)

	if (cond1 || cond2) && (structure.Port == structure.QlitePort) {
		err = errors.New("qlite port and enqueuer port are the same")
		return
	}

	address = fmt.Sprintf("%s:%d", localIp, structure.Port)

	driver, err = qldriver.NewPublisherClient(
		structure.Name,
		structure.Password,
		structure.MaxConnections,
		structure.MaxIoTimeSeconds,
		structure.QlitePort,
		structure.QliteAddress)

	return
}

func main() {

	fPath := os.Args[1]
	data, err := os.ReadFile(fPath)

	if err != nil {
		log.Fatal("could not read configuration file")
	}

	qDriver, key, address, err := readConfig(data)

	if err != nil {
		log.Fatal(err)
	}

	handler := enqueueHandler{
		key:    key,
		driver: qDriver,
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handler)

	err = http.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}
}
