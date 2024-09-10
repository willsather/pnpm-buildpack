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
    pack build sample-app -p ./pnpm-install/integration/sample-app -b gcr.io/paketo-buildpacks/node-engine -b ./pnpm/build/pnpm-buildpack.cnb -b ./pnpm-install/build/pnpm-install-buildpack.cnb -b ./pnpm-start/build/pnpm-start-buildpack.cnb
    ```


## TODO
- [ ] add/update `pnpm-install` unit tests
- [ ] add/update `pnpm-start` unit tests
- [ ] fix `pnpm-install` integration test(s)
- [ ] add `pnpm-start` integration test(s)
- [ ] if, how, and where to add `pnpm build` step??
- [ ] fix `build` vs `launch` dependencies
- [ ] get `pnpm` version from `package.json` in `pnpm-install/detect.go`
- [ ] add latest `pnpm` dependency
- [ ] create primary script to bundle all three buildpacks
- [ ] documents steps for building / packaging / publishing buildpack
- [ ] add more integration 
