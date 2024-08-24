# PNPM Buildpack

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