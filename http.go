package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed webConfig.html
var webConfig []byte

func NewHTTPServer() *http.Server {
	mux := mux.NewRouter()

	mux.HandleFunc("/", handleHome).
		Methods(http.MethodGet)
	mux.HandleFunc("/switches", handleGet).
		Methods(http.MethodGet)
	mux.HandleFunc("/switches", handlePut).
		Methods(http.MethodPut)

	return &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write(webConfig)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	b, err := ReadConfig("configuration.yaml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := Conf{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	conf := Conf{}
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := ReadConfig("configuration.yaml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yamlMap := yaml.MapSlice{}
	yaml.Unmarshal(b, &yamlMap)

	switchIndex := -1
	for i := range yamlMap {
		if yamlMap[i].Key == "switch" {
			switchIndex = i
			break
		}
	}

	sws := []yaml.MapSlice{}
	for j := range conf.Switch {
		sws = append(sws, conf.Switch[j].ToMapSlice())
	}

	if switchIndex >= 0 {
		yamlMap[switchIndex].Value = sws
	} else {
		sw := yaml.MapItem{Key: "switch", Value: sws}
		yamlMap = append(yamlMap, sw)
	}

	err = WriteConfig("configuration.yaml", yamlMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ActivateConfig("configuration.yaml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type Conf struct {
	Switch []Switch `yaml:"switch" json:"switch"`
}

type Switch struct {
	Platform              string `yaml:"platform" json:"platform"`
	Name                  string `yaml:"name" json:"name"`
	Unique_id             string `yaml:"unique_id" json:"unique_id"`
	Command_topic         string `yaml:"command_topic" json:"command_topic"`
	State_topic           string `yaml:"state_topic" json:"state_topic"`
	Availability_topic    string `yaml:"availability_topic" json:"availability_topic"`
	Payload_on            string `yaml:"payload_on" json:"payload_on"`
	Payload_off           string `yaml:"payload_off" json:"payload_off"`
	Payload_available     string `yaml:"payload_available" json:"payload_available"`
	Payload_not_available string `yaml:"payload_not_available" json:"payload_not_available"`
	Qos                   int    `yaml:"qos" json:"qos"`
	Retain                bool   `yaml:"retain" json:"retain"`
}

func (s *Switch) ToMapSlice() yaml.MapSlice {
	return yaml.MapSlice{
		{Key: "platform", Value: s.Platform},
		{Key: "name", Value: s.Name},
		{Key: "unique_id", Value: s.Unique_id},
		{Key: "command_topic", Value: s.Command_topic},
		{Key: "state_topic", Value: s.State_topic},
		{Key: "availability_topic", Value: s.Availability_topic},
		{Key: "payload_on", Value: s.Payload_on},
		{Key: "payload_off", Value: s.Payload_off},
		{Key: "payload_available", Value: s.Payload_available},
		{Key: "payload_not_available", Value: s.Payload_not_available},
		{Key: "qos", Value: s.Qos},
		{Key: "retain", Value: s.Retain},
	}
}
