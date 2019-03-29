package got

var CsetPaths = (*Config).setPaths
var CaddPath = (*Config).addPath
var Cload = (*Config).load

func (c *Config) ExportCusedPath() string {
	return c.usedPath
}

func (c *Config) ExportCpaths() []string {
	return c.paths
}
