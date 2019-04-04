package got

type Manager interface {
	Install(name string) error
}

type apt struct {
}

func (a *apt) Install(name string) error {
	c := NewCommand()
	return c.SURun("apt-get", "install", name)
}

func NewManager(manager string) Manager {
	switch manager {
	case "apt", "apt-get":
		return &apt{}
	default:
		return nil
	}
}
