package volume

type Config struct {
	His int
}

func (c Config) Ensure() Config {
	if c.His == 0 {
		c.His = 10
	}

	return c
}

func (c Config) Verify() {
	if c.His == 0 {
		panic("Config.His must not be empty")
	}
}
