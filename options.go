package configuration

// ConfiguratorOption defines Option function for Configuration
type ConfiguratorOption func(*Configurator)

// OnFailFnOpt sets function which will be called when an error occurs during Configurator.applyProviders()
func OnFailFnOpt(fn func(error)) ConfiguratorOption {
	return func(c *Configurator) {
		c.onErrorFn = fn
	}
}
