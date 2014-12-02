package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chef/gohai/cpu"
	"github.com/go-chef/gohai/filesystem"
	"github.com/go-chef/gohai/kernel"
	"github.com/go-chef/gohai/memory"
	"github.com/go-chef/gohai/network"
	"github.com/go-chef/gohai/password"
	"github.com/go-chef/gohai/platform"
	"github.com/go-chef/gohai/util"
	"log"
	"os"
	"os/exec"
)

// An interface all gohai submodules must satisfy.
type collector interface {
	Name() string
	Collect() (interface{}, error)
	Provides() []string
}

// The built-in gohai information collectors to run.
var collectors = []collector{
	&cpu.Cpu{},
	&filesystem.FileSystem{},
	&memory.Memory{},
	&network.Network{},
	&kernel.Kernel{},
	&password.Password{},
	&platform.Platform{},
}

// The default plugin directory. TODO: make configurable.
var defaultPluginDir = "/var/lib/gohai/plugins"

// Run the data collectors from above, merge their information, and return it
// to the caller.
func collect() (map[string]interface{}, error) {
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
				err := util.MergeMap(result, c)
				if err != nil {
					return nil, err
				}
			default: // try?
				result[collector.Name()] = c
			}
		}
	}

	return result, nil
}

func main() {
	// the built-in data colletors are run here
	gohai, err := collect()

	if err != nil {
		panic(err)
	}

	// now running the external plugins
	err = runPlugins(gohai)
	if err != nil {
		panic(err)
	}

	// marshal to json, print, and exit
	buf, err := json.Marshal(gohai)

	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	json.Indent(&out, buf, "", "  ")

	out.WriteTo(os.Stdout)
	os.Exit(0)
}

// Run the external plugins found in the plugin dir. These plugins need to
// return a JSON hash of data.
func runPlugins(gohai map[string]interface{}) error {
	// TODO: make plugin dir configurable
	pDir, err := os.Open(defaultPluginDir)
	if err != nil {
		// Only return an error if Open returned an error that wasn't
		// "directory does not exist". If the plugin dir doesn't exist,
		// just return nil, but obviously don't go trying to run the
		// plugins.
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	// get the names of the plugins to run
	pRun, err := pDir.Readdirnames(0)
	if err != nil {
		return err
	}
	pRunLen := len(pRun)
	if pRunLen == 0 {
		// no plugins, bail now
		return nil
	}

	outCh := make(chan map[string]interface{}, pRunLen)

	// For each plugin, spawn a goroutine to run the plugin and place its
	// returned data into the channel.
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
			outCh <- jsonOut
		}()
	}
	// Process the returned data from the plugins as it comes in.
	for i := 0; i < pRunLen; i++ {
		out := <-outCh
		err := util.MergeMap(gohai, out)
		if err != nil {
			return nil
		}
	}

	return nil
}
