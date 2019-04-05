package got

// Manager defines the operation of package manager.
type Manager interface {
	Install(names ...string) error
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

func (a *apt) Install(names ...string) error {
	c := NewCommand()
	args := make([]string, len(names)+1)
	args[0] = "install"
	args = append(args, names...)
	return c.SURun("apt-get", args...)
}
