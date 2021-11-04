package configuration

// ConfiguratorOption defines Option function for Configuration
type ConfiguratorOption func(*Configurator)

// LoggerOpt sets a custom logger
func LoggerOpt(l func(format string, v ...interface{})) ConfiguratorOption {
	return func(c *Configurator) {
		c.loggerFn = l
	}
}

// EnableLoggingOpt enables/disables logging
func EnableLoggingOpt(enable bool) ConfiguratorOption {
	return func(c *Configurator) {
		c.loggingEnabled = enable
	}
}

// OnFailFnOpt sets function which will be called when an error occurs during Configurator.applyProviders()
func OnFailFnOpt(fn func(error)) ConfiguratorOption {
	return func(c *Configurator) {
		c.onErrorFn = fn
	}
}
