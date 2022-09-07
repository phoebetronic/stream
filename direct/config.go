package direct

import "time"

type Config struct {
	Dur time.Duration
	His int
}

func (c Config) Ensure() Config {
	if c.Dur == 0 {
		c.Dur = time.Minute
	}

	if c.His == 0 {
		c.His = 5
	}

	return c
}

func (c Config) Verify() {
	if c.Dur == 0 {
		panic("Config.Dur must not be empty")
	}
	if c.His == 0 {
		panic("Config.His must not be empty")
	}
}
