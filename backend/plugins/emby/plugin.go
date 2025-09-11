package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"emby-plugin/internal/plugin"
	"emby-plugin/internal/pluginsdk"

	"github.com/mitchellh/mapstructure"
)

func NewPlugin() (pluginsdk.Plugin, error) {
	return &plugin.EmbyPlugin{}, nil
}

func main() { test() }

func test() {
	p, err := NewPlugin()
	if err != nil {
		panic(err)
	}
	dataSource, err := os.ReadFile("./tv.json")
	if err != nil {
		panic(err)
	}
	var data map[string]any
	if err = json.Unmarshal(dataSource, &data); err != nil {
		panic(err)
	}
	settingData, err := os.ReadFile("./setting.json")
	if err != nil {
		panic(err)
	}
	var config map[string]any
	if err = json.Unmarshal(settingData, &config); err != nil {
		panic(err)
	}
	var settings map[string]any
	if err = mapstructure.Decode(config["settings"], &settings); err != nil {
		panic(err)
	}
	res, err := p.Process(context.Background(), data, settings)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------------------------------")
	fmt.Println(res.Title)
	fmt.Println(res.Content)
	fmt.Println(res.Image)
	fmt.Println(res.URL)
	fmt.Println(res.Targets)
}
