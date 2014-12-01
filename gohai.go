package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/go-chef/gohai/cpu"
	"github.com/go-chef/gohai/filesystem"
	"github.com/go-chef/gohai/kernel"
	"github.com/go-chef/gohai/memory"
	"github.com/go-chef/gohai/network"
	"github.com/go-chef/gohai/password"
	"github.com/go-chef/gohai/platform"
)

type Collector interface {
	Name() string
	Collect() (interface{}, error)
	Provides() []string
}

var collectors = []Collector{
	&cpu.Cpu{},
	&filesystem.FileSystem{},
	&memory.Memory{},
	&network.Network{},
	&kernel.Kernel{},
	&password.Password{},
	&platform.Platform{},
}

var defaultPluginDir = "/var/lib/gohai/plugins"

func Collect() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, collector := range collectors {
		c, err := collector.Collect()
		if err != nil {
			log.Printf("[%s] %s", collector.Name(), err)
			continue
		}
		// password stuff is nil in windows, I believe
		if c != nil {
			switch c := c.(type) {
			case map[string]interface{}:
				// TODO: needs a real merge function
				for k, v := range c {
					result[k] = v
				}
			default: // try?
				result[collector.Name()] = c
			}
		}
	}

	return result, nil
}

func main() {
	gohai, err := Collect()

	if err != nil {
		panic(err)
	}

	err = runPlugins(gohai)
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

func runPlugins(gohai map[string]interface{}) error {
	// TODO: make plugin dir configurable
	pDir, err := os.Open(defaultPluginDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	pRun, err := pDir.Readdirnames(0)
	if err != nil {
		return err
	}

	outCh := make(chan map[string]interface{}, len(pRun))

	for _, v := range pRun {
		go func() {
			cmdStr := defaultPluginDir + "/" + v
			cmd := exec.Command(cmdStr)

			output, err := cmd.Output()
			if err != nil {
				log.Println(err)
				return
			}
			jsonOut := make(map[string]interface{})
			err = json.Unmarshal(output, &jsonOut)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("sending output...")
			outCh <- jsonOut
		}()
	}
	for i := 0; i < len(pRun); i++ {
		out := <- outCh
		log.Println("received output")
		// TODO: needs real merge of course
		for k, v := range out {
			gohai[k] = v
		}
	}
	log.Println("done with plugins")

	return nil
}
