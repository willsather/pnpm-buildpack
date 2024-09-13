# PNPM Buildpack

This repository is a monorepo of different modules needed to install necessary dependencies and run Node application
that use `pnpm` as a package manager.

## Modules

This repository contains three different Cloud Native Buildpacks:

1. `pnpm`
2. `pnpm-install`
3. `pnpm-start`

Which are can be built separately, but can be combined to fully build a Node application that utilizes the `pnpm`
package manager.

Roughly:

1. `pnpm` provides the `pnpm` dependency as a global dependency on the image.

2. `pnpm-install` provides `node_modules` and downloads the necessary dependencies on the image.

3. `pnpm-start` takes the image that has `pnpm` and `node_modules` and runs the `start` command for the Node
   application.

## Local Setup

Current testing looks like:

1. Build each buildpack using (_this can also be done manually within each buildpack's `./scripts/build.sh`_)

   ```bash
   ./scripts/build.sh
   ```

2. Package `pnpm-buildpack` (_composite buildpack using `pnpm`, `pnpm-install`, and `pnpm-start`_)

   ```bash
   ./scripts/package.sh
   ```

3. Create image for sample application:
    ```bash
    pack build pnpm-simple-app -p ./integration/simple-app -b ./build/pnpm-buildpack.cnb
    ```

where `pnpm-simple-app` is the image name, `./integration/simple-app` is the location of the project, and
`./build/pnpm-buildpack.cnb` is the output of the combined `pnpm-buildpack` from step 2.

## TODO

### Fix / Cleanup

- [x] add `node_mdules/.bin` to path (how does this work with `pnpm`?)
- [ ] fix `build` vs `launch` dependencies
- [x] fix what happens in `pnpm-install` if `node_modules` already exists (_what happens if it exists but is invalid?_)
- [ ] add latest `pnpm` dependency

### Testing

- [ ] fix `pnpm-install` integration test(s)
- [ ] add `pnpm-start` integration test(s)
- [ ] add more full integration projects

### Features

- [ ] cache `pnpm install` on layer (if the `sha` is the same)
- [ ] generate sbom for dependencies / CVEs / auditability

### Building / Publishing

- [x] create primary script to bundle all three buildpacks
- [x] documents steps for building / packaging / publishing buildpack

## Questions

- [x] how are workspaces handled in existing buildpacks (npm/yarn)
- [x] how do buildpacks get published? which registry? (how does `cf create-buildpack` work in this context)
- [x] how do multiple buildpacks get combined into just a single one like (`nodejs_buildpack`)
