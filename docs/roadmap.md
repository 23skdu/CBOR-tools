# Roadmap

This roadmap tracks remaining work for CBOR-tools after the initial Go rewrite.

## Stage 1: Developer experience (done)

- [x] Go module and CLI skeleton
- [x] Unit and fuzz tests at 100% coverage
- [x] Makefile targets for build/test/fuzz/cover/clean/lint/security
- [x] gofmt, vet, and tidy checks
- [x] CLI version flag
- [x] Input size guard

## Stage 2: Security hardening

- [ ] Pin GitHub Actions to commit SHAs
- [ ] Add Trivy container/filesystem scan workflow
- [ ] Generate and publish SBOM (Syft/CycloneDX)
- [ ] Add CodeQL analysis workflow
- [ ] Harden file access with `os.Root` where practical
- [ ] Add dependency review workflow

## Stage 3: Release automation

- [ ] Add GoReleaser configuration
- [ ] Add release workflow on tags
- [ ] Produce cross-platform binaries
- [ ] Publish Docker images to registry
- [ ] Generate checksums and signatures

## Stage 4: CLI enhancements

- [ ] Add `--format` output option (json, pretty, raw)
- [ ] Add atomic output file writes
- [ ] Add shell completion generation
- [ ] Add man page generation
- [ ] Add `--quiet` mode

## Stage 5: Testing and quality

- [ ] Add benchmark tests
- [ ] Add example tests
- [ ] Add CLI-level fuzz targets
- [ ] Add golden file tests
- [ ] Add integration tests for Docker image

## Stage 6: Documentation and community

- [ ] Add man page
- [ ] Add usage examples directory
- [ ] Add issue and PR templates
- [ ] Add CODEOWNERS
