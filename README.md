# PNPM Buildpack

This repository is a monorepo of different modules needed to install necessary dependencies and run Node application
that use `pnpm` as a package manager.

> [!WARNING]  
> This repository is in active development, there is no expectation this code is functioning or usable (yet). :)

## Modules

This repository contains three different Cloud Native Buildpacks:

1. `pnpm`
2. `pnpm-install`
3. `pnpm-start`

Which are can be built separately, but can be combined to fully build a Node application that utilizes the `pnpm`
package manager.

Roughly, `pnpm` provides the `pnpm` dependency as a global dependency on the image.  `pnpm-install` provides
`node_modules` and downloads the necessary dependencies on the image. And lastly, `pnpm-start` takes the image that has
`pnpm` and `node_modules` and runs the `start` command for the Node application. 

## Local Setup

Current testing looks like:

1. build and package each module

2. pack sample application:
    ```bash
    pack build simple-app -p ./integration/simple-app -b gcr.io/paketo-buildpacks/node-engine -b ./pnpm/build/pnpm-buildpack.cnb -b ./pnpm-install/build/pnpm-install-buildpack.cnb -b ./pnpm-start/build/pnpm-start-buildpack.cnb
    ```


## TODO

### Cleanup
- [ ] if, how, and where to add `pnpm build` step??
- [ ] fix what happens in `pnpm-install` if `node_modules` already exists (_what happens if it exists but is invalid?_)
- [ ] fix `build` vs `launch` dependencies
- [ ] add latest `pnpm` dependency

### Testing
- [ ] fix `pnpm-install` integration test(s)
- [ ] add `pnpm-start` integration test(s)
- [ ] add more full integration projects

### Features
- [ ] cache `pnpm install` on layer (if the sha is the same)

### Building / Publishing
- [ ] create primary script to bundle all three buildpacks
- [ ] documents steps for building / packaging / publishing buildpack
- [ ] better utilize `jam` CLI tool to build/package each buildpack


## Questions
- [ ] how are workspaces handled in existing buildpacks (npm/yarn)
- [ ] how do buildpacks get published? which registry? (how does `cf create-buildpack` work in this context)
- [ ] how do multiple buildpacks get combined into just a single one like (`nodejs_buildpack`)
