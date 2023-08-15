WIP  TBD

 * Adding zedpm-plugin-golanci to run the golangci-lint command.

v0.1.1  2023-08-15

 * Fix default configuration to look for zedpm-plugin-go

v0.1.0  2023-08-15

 * The zedpm command handles quit via signal better.
 * Added zedpm run test command via plugin-go.
 * Improved error messages in plugin-git.
 * Added a better UI for tracking progress of a goal.
 * Added the --log-level option to all commands to adjust logging level.
 * Added the --log-file option to enable raw logs to file.
 * Added the --progress option to let the end-user disable the new UI and see 
   raw logs instead.
 * Fix: zedpm run release now completes the final release tag push properly.

v0.0.0  2023-03-29

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
