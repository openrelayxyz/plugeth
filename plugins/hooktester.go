package plugins


// type PluginLoader struct{
// 	Plugins []*plugin.Plugin
// 	Subcommands map[string]Subcommand
// 	Flags []*flag.FlagSet
// 	LookupCache map[string][]interface{}
// }

func HookTester(name string, fn interface{}) func() {
  oldDefault := DefaultPluginLoader
  DefaultPluginLoader = &PluginLoader{
    LookupCache: map[string][]interface{}{
      name: []interface{}{fn},
    },
  }
  return func() { DefaultPluginLoader = oldDefault }
}
