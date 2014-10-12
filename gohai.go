package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/ctdk/gohai/cpu"
	"github.com/ctdk/gohai/filesystem"
	"github.com/ctdk/gohai/kernel"
	"github.com/ctdk/gohai/memory"
	"github.com/ctdk/gohai/network"
	"github.com/ctdk/gohai/platform"
)

type Collector interface {
	Name() string
	Collect() (interface{}, error)
}

var collectors = []Collector{
	&cpu.Cpu{},
	&filesystem.FileSystem{},
	&memory.Memory{},
	&network.Network{},
	&kernel.Kernel{},
}

var topLevelCollectors = []Collector{
	&platform.Platform{},
}

func Collect() (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	for _, collector := range collectors {
		c, err := collector.Collect()
		if err != nil {
			log.Printf("[%s] %s", collector.Name(), err)
			continue
		}
		result[collector.Name()] = c
	}
	// platform is weird, this stuff is top level
	for _, collector := range topLevelCollectors {
		c, err := collector.Collect()
		if err != nil {
			log.Printf("[%s] %s", collector.Name(), err)
			continue
		}
		switch c := c.(type) {
		case map[string]interface{}:
			for k, v := range c {
				result[k] = v
			}
		default:
			result[collector.Name()] = c
		}
	}

	return
}

func main() {
	gohai, err := Collect()

	if err != nil {
		panic(err)
	}

	buf, err := json.Marshal(gohai)

	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	json.Indent(&out, buf, "", "  ")

	out.WriteTo(os.Stdout)
}
