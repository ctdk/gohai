package cpu

type Cpu struct{}

const name = "cpu"

func (self *Cpu) Name() string {
	return name
}

func (self *Cpu) Collect() (interface{}, error) {
	result, err := getCpuInfo()

	return result, err
}

func (self *Cpu) Provides() []string {
	return []string{"cpu"}
}
