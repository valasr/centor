package pluginManager

import (
	apiCaller_plugin "github.com/mrtdeh/centor/plugins/apiCaller"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
	installer_plugin "github.com/mrtdeh/centor/plugins/installer"
	system_plugin "github.com/mrtdeh/centor/plugins/system"
	timeSyncer_plugin "github.com/mrtdeh/centor/plugins/time_syncer"
)

type Config struct {
	PluginKits.Config
}

func Bootstrap(cnf Config) error {
	return PluginKits.Loader(cnf.Config, func(pms *PluginKits.PluginManagerService) {

		pms.AddPlugin(&installer_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "installer",
			},
		})

		pms.AddPlugin(&timeSyncer_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "time-syncer",
			},
		})

		pms.AddPlugin(&system_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "system",
			},
		})

		pms.AddPlugin(&apiCaller_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "api-caller",
			},
		})

	})
}
