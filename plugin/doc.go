// Package plugin is the container for all things plugins. The top level package
// defines the interface for plugins. A plugin must, at the very least, provide
// an implementation of plugin.Interface. Then, that plugin is run by providing
// a main function that runs the plugin using the RunPlugin() function defined
// in "github.com/zostay/zedpm/plugin/metal".
//
// From there, the plugin may:
//
// * Define task implementations via the Implements, Prepare, Cancel, and
// Complete Prepare methods of the plugin.Interface.
//
// * Define goals via the Goal method.
//
// Plugins may be written in other languages than Golang. To do that, the plugin
// will need to implement the plugin interface defined by
// "github.com/hashicorp/go-plugin" to handle the startup of a gRPC server and
// then printing details regarding that gRPC server on stdout. The gRPC server
// implemented must implement the interface defined in the TaskExecution service
// found in task-interface.proto.
package plugin

// TODO There should be an example plugin written in python somewhere to show this is possible, prove it works, and provide an example how to do it.
