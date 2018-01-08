package tftest

type Tester interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type TestContext struct {
	*Context
	t Tester
}

func NewTestContext(t Tester, options ...ContextOption) *TestContext {
	return &TestContext{
		Context: NewContext(options...),
		t:       t,
	}
}

func (c *TestContext) Apply() {
	if _, err := c.Context.Apply(); err != nil {
		c.t.Fatal(err)
	}
}

func (c *TestContext) Destroy() {
	if _, err := c.Context.Destroy(); err != nil {
		c.t.Fatal(err)
	}
}

func (c *TestContext) Import(resource, id string) {
	if _, err := c.Context.Import(resource, id); err != nil {
		c.t.Fatal(err)
	}
}

func (c *TestContext) Output(name string) string {
	output, err := c.Context.Output(name)
	if err != nil {
		c.t.Fatal(err)
	}

	return output
}
