package main

import (
	"io/fs"
	"io/ioutil"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func ReadConfig(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func WriteConfig(path string, m yaml.MapSlice) error {
	o, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, o, fs.ModePerm)
}

func ActivateConfig(path string) error {
	err := exec.Command("/bin/bash", "-c", "docker cp configuration.yaml homeassistant:config/configuration.yaml").Run()
	if err != nil {
		return err
	}

	err = exec.Command("/bin/bash", "-c", "docker restart homeassistant").Run()
	if err != nil {
		return err
	}

	return nil
}

func FetchConfig(path string) error {
	err := exec.Command("/bin/bash", "-c", "docker cp homeassistant:config/configuration.yaml "+path).Run()
	if err != nil {
		return err
	}

	return nil
}
