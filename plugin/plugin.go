package plugin

type Info struct{}

var InfoCh chan map[string]interface{}
var DefaultPluginDir = "/var/lib/gohai/plugins"
var DefaultSocket = "/tmp/gohai.sock"
var DefaultAddr = "127.0.0.1:9966"

func init() {
	InfoCh = make(chan map[string]interface{}, 1)
}

func (i *Info) SendInfo(args map[string]interface{}, reply *string) error {
	InfoCh <- args
	*reply = "ok"
	return nil
}
