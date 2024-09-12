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

1. Build each buildpack (manually or use `./scripts/build.sh`)

2. Package `pnpm-buildpack` (combination of four buildpacks)

3. Create image for sample application:
    ```bash
    pack build pnpm-simple-app -p ./integration/simple-app -b ./build/pnpm-buildpack.cnb
    ```

where `pnpm-simple-app` is the image name, `./integration/simple-app` is the location of the project, and
`./build/pnpm-buildpack.cnb` is the output of the combined `pnpm-buildpack` from step 2. 

## TODO

### Cleanup

- [ ] fix what happens in `pnpm-install` if `node_modules` already exists (_what happens if it exists but is invalid?_)
- [ ] fix `build` vs `launch` dependencies
- [ ] add latest `pnpm` dependency

### Testing

- [ ] fix `pnpm-install` integration test(s)
- [ ] add `pnpm-start` integration test(s)
- [ ] add more full integration projects

### Features

- [ ] cache `pnpm install` on layer (if the sha is the same)
- [ ] generate sbom for dependencies / CVEs / auditability

### Building / Publishing

- [x] create primary script to bundle all three buildpacks
- [x] documents steps for building / packaging / publishing buildpack
- [ ] better utilize `jam` CLI tool to build/package each buildpack

## Questions

- [ ] how are workspaces handled in existing buildpacks (npm/yarn)
- [ ] how do buildpacks get published? which registry? (how does `cf create-buildpack` work in this context)
- [ ] how do multiple buildpacks get combined into just a single one like (`nodejs_buildpack`)
