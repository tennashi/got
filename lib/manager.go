package got

// Manager defines the operation of package manager.
type Manager interface {
	Install(name string) error
}

// NewManager initialize Manager from the command name.
func NewManager(manager string) Manager {
	switch manager {
	case "apt", "apt-get":
		return &apt{}
	default:
		return nil
	}
}

type apt struct {
}

func (a *apt) Install(name string) error {
	c := NewCommand()
	return c.SURun("apt-get", "install", name)
}
