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

- [x] Pin GitHub Actions to commit SHAs
- [x] Add Trivy container scan workflow
- [ ] Generate and publish SBOM (Syft/CycloneDX)
- [x] Add CodeQL analysis workflow
- [ ] Harden file access with `os.Root` where practical
- [x] Add dependency review workflow

## Stage 3: Release automation

- [x] Add GoReleaser configuration
- [x] Add release workflow on tags
- [x] Produce cross-platform binaries
- [ ] Publish Docker images to registry
- [ ] Generate checksums and signatures

## Stage 4: CLI enhancements

- [x] Add `--format` output option (pretty, json)
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
- [x] Add issue and PR templates
- [x] Add CODEOWNERS