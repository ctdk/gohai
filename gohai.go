package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/go-chef/gohai/cpu"
	"github.com/go-chef/gohai/filesystem"
	"github.com/go-chef/gohai/kernel"
	"github.com/go-chef/gohai/memory"
	"github.com/go-chef/gohai/network"
	"github.com/go-chef/gohai/password"
	"github.com/go-chef/gohai/platform"
	"github.com/go-chef/gohai/plugin"
)

type Collector interface {
	Name() string
	Collect() (interface{}, error)
	Provides () []string
}

var collectors = []Collector{
	&cpu.Cpu{},
	&filesystem.FileSystem{},
	&memory.Memory{},
	&network.Network{},
	&network.Counters{},
	&kernel.Kernel{},
	&password.Password{},
	&platform.Platform{},
}

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
	pDir, err := os.Open(plugin.DefaultPluginDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	pRun, err := pDir.Readdirnames(0)
	if err != nil {
		return err
	}

	stopch := make(chan struct{}, 1)
	donech := make(chan struct{}, 1)
	readych := make(chan struct{}, 1)
	go startPluginServer(stopch, readych, donech)

	go func() {
		for {
			pi := <-plugin.InfoCh
			// TODO: make a merge function. Also, a mutex for the
			// info hash.
			for k, v := range pi {
				gohai[k] = v
			}
		}
	}()
	<-readych

	for _, v := range pRun {
		cmdStr := plugin.DefaultPluginDir + "/" + v
		// TODO: make socket/network addr passable to plugin
		cmd := exec.Command(cmdStr)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log.Println(stderr.String())
			return err
		}
	}

	stopch <- struct{}{}
	<-donech
	return nil
}

func startPluginServer(stopch <-chan struct{}, readych, donech chan<- struct{}) {
	// TODO: make socket/network addr configurable
	uaddr, _ := net.ResolveUnixAddr("unix", plugin.DefaultSocket)
	l, err := net.ListenUnix("unix", uaddr)
	readych <- struct{}{}
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		<-c
		l.Close()
		os.Exit(0)
	}(sigch)

	if err != nil {
		log.Printf("Failed to start socket for plugins: %s\n", err.Error())
		os.Exit(1)
	}
	rpc.Register(new(plugin.Info))
	done := false
	go func() {
		for {
			log.Printf("Waiting for plugins...")
			if conn, err := l.AcceptUnix(); err == nil {
				log.Println("reading data from plugin...")
				go jsonrpc.ServeConn(conn)
			} else {
				if !done {
					log.Printf("Plugin connection failed: %s ", err.Error())
					os.Exit(1)
				} else {
					return
				}
			}

		}
	}()
	<-stopch
	done = true
	err = l.Close()
	donech <- struct{}{}
	if err != nil {
		log.Println("err closing ", err)
	}
	return
}
