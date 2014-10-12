package plugin

type Info struct {}

var InfoCh chan map[string]interface{}
var PluginDir = "/tmp/plugins"

func init() {
	InfoCh = make(chan map[string]interface{}, 1)
}

func (i *Info) SendInfo(args map[string]interface{}, reply *string) error {
	InfoCh <- args
	*reply = "ok"
	return nil
}
