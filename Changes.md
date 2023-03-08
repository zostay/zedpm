WIP  TBD

 * Initial release.
 * Provides the zedpm command entrypoint with self-configuring sub-commands
   based on configured plugins.
 * Defines the gRPC protocol used to communicate between the zedpm command and
   plugins, defined according to Hashicorp's go-plugin system.
 * Defines tooling to make writing Golang plugins simple and make plugins in
   other languages possible.
 * Tools for loading HCL formatted configuration files to hierarchically define
   configuration for goals, tasks, targets, and plugins.
 * Defines the zedpm-plugin-changelog plugin.
 * Defines the zedpm-plugin-git plugin.
 * Defines the zedpm-plugin-github plugin.
 * Defines the zedpm-plugin-goals plugin.
 * Defines the following goals: build, deploy, generate, info, init, install,
   lint, release, request, and test.
 * Implements parts of the info and lint goals and a reasonably complete release
   goal.
 * Provides documentation and design drawings to point to the future.
 * And so much more!
