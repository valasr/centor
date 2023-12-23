package PluginKits

func Loader(cnf Config, load func(*PluginManagerService)) error {
	if err := validateConfig(cnf); err != nil {
		return err
	}
	// Create a new instance of pman
	pman := &PluginManagerService{}

	// Add plugins to the pman
	load(pman)

	// Start the pman and its plugins
	pman.Start(cnf.CoreHandler, cnf.RouterAPI)
	return nil
}
