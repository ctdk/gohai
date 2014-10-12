package kernel

type Kernel struct{}

const name = "kernel"

func (k *Kernel) Name() string {
	return name
}

func (k *Kernel) Collect() (interface{}, error) {
	result, err := getKernelInfo()
	return result, err
}
