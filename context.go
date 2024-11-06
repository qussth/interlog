package interlog

// With func
// Create new context with values, which will be appended to log message and remains after call of log method
func (l *Logger) With(values []Value) *Context {
	return &Context{
		l:      l,
		values: values,
	}
}

func (c *Context) Append(values []Value) *Context {
	c.values = append(c.values, values...)

	return c
}

func (c *Context) Flush() {
	c.values = c.values[:0:0]
}

func (c *Context) Debug(message string, values []Value) {
	c.l.Debug(message, append(c.values, values...))
}

func (c *Context) Info(message string, values []Value) {
	c.l.Info(message, append(c.values, values...))
}

func (c *Context) InfoToSentry(message string, values []Value) {
	c.l.InfoToSentry(message, append(c.values, values...))
}

func (c *Context) Warn(message string, values []Value) {
	c.l.Warn(message, append(c.values, values...))
}

func (c *Context) Error(err error, values []Value) {
	c.l.Error(err, append(c.values, values...))
}

func (c *Context) Panic(err error, values []Value) {
	c.l.Panic(err, append(c.values, values...))
}
