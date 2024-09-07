# PNPM Install

This PNPM Install Cloud Native Buildpack is responsible for generating and providing application dependencies for Node
applications that use the `pnpm` package manager.

## Integration

The PNPM Install Cloud Native Buildpack provides `node_modules` as a dependency. Downstream buildpacks can require the
`node_modules` dependency using a Build Plan TOML file that looks like the following:

```toml
[[requires]]
   # The name of the PNPM Install dependency is "node_modules".
   name = "node_modules"
   
   # The PNPM Install buildpack supports some non-required metadata options.
   [requires.metadata]
      # Setting the build flag to true will ensure that the node modules
      # are available for subsequent buildpacks during their build phase.
      # If you are writing a buildpack that needs a node module during
      # its build process, this flag should be set to true.
      build = true
      
      # Setting the launch flag to true will ensure that the packages
      # managed by PNPM are available for the running application. If you
      # are writing an application that needs node modules at runtime,
      # this flag should be set to true.
      launch = true
```

Downstream buildpacks can also utilize this Cloud Native Buildpack by adding this Buildpack Plan as PlanRequirement like
this:

```go
package pnpm

import (
	"github.com/paketo-buildpacks/packit/v2"
)

func Detect() packit.DetectFunc {
   return func(context packit.DetectContext) (packit.DetectResult, error) {
      return packit.DetectResult{
         Plan: packit.BuildPlan{
            Requires: []packit.BuildPlanRequirement{
               {
                  Name: "node_modules", // HERE!
                  Metadata: map[string]interface{}{
                     "build": true,
                     "launch": true,
                  },
               },
            },
         },
      }, nil
   }
}
```

## Local Setup

1. Install [Pack CLI](https://buildpacks.io/docs/for-platform-operators/how-to/integrate-ci/pack/)
2. Set Default Builder

   ```zsh
   pack config default-builder paketobuildpacks/builder-jammy-base
   ```

3. Verify Pack CLI Default Builder
   ```zsh
   pack config default-builder
   ```

4. Build Go Application
   ```zsh
   ./scripts/build.sh
   ```

5. Package Buildpack
   ```zsh
   ./scripts/package.sh
   ```