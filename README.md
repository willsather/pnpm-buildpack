# PNPM Buildpack

This repository is a monorepo of different modules needed to install necessary dependencies and run Node application
that use `pnpm` as a package manager.

> [!WARNING]  
> This repository is in active development, there is no expectation this code is functioning or usable (yet). :)

## Local Setup

TODO: create scripts and steps for building / packaging / publishing buildpack

Current testing looks like:

1. build and package each module

2. pack sample application:
    ```bash
    pack build sample-app -p ./pnpm-install/integration/sample-app -b gcr.io/paketo-buildpacks/node-engine -b ./pnpm/build/pnpm-buildpack.cnb -b ./pnpm-install/build/pnpm-install-buildpack.cnb -b ./pnpm-start/build/pnpm-start-buildpack.cnb
    ```