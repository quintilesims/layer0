package tftest

type ContextOption func(*Context)

func Log(logger Logger) ContextOption {
	return func(c *Context) {
		c.Logger = logger
	}
}

func Dir(dir string) ContextOption {
	return func(c *Context) {
		c.dir = dir
	}
}

func DryRun(dryRun bool) ContextOption {
	return func(c *Context) {
		c.dryRun = dryRun
	}
}

func Var(name, val string) ContextOption {
	return func(c *Context) {
		c.Vars[name] = val
	}
}

func Vars(vars map[string]string) ContextOption {
	return func(c *Context) {
		for name, val := range vars {
			c.Vars[name] = val
		}
	}
}
