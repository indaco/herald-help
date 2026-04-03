# Variables

export LOG_STYLE := "emoji"
logger := "scripts/lib/logger.sh"

# Go commands

go := "go"
goclean := go + " clean"

# Default - show help
default:
    @just --list

# Clean the build directory and Go cache
clean:
    @. {{ logger }} && log_info "Clean the build directory and Go cache"
    rm -f coverage.out coverage.html
    {{ goclean }} -cache

# === Code Quality ===

# Format code
fmt:
    @. {{ logger }} && log_info "Running fmt"
    {{ go }} fmt ./...

# Run go-modernize with auto-fix
modernize:
    @. {{ logger }} && log_info "Running go-modernize"
    modernize --fix ./...

# Run golangci-lint
lint:
    @. {{ logger }} && log_info "Running golangci-lint"
    golangci-lint run ./...

# Run goreportcard-cli
reportcard:
    @. {{ logger }} && log_info "Running goreportcard-cli"
    goreportcard-cli -v

# Run govulncheck
security-scan:
    @. {{ logger }} && log_info "Running govulncheck"
    govulncheck ./...

# Run modernize, lint, and reportcard
check: fmt modernize lint reportcard

# Run go mod tidy
tidy:
    @. {{ logger }} && log_info "Running go mod tidy"
    {{ go }} mod tidy

# Run go mod download
deps:
    @. {{ logger }} && log_info "Running go mod download"
    {{ go }} mod download

# === Test Recipes ===

# Run all tests and print code coverage value
test:
    @. {{ logger }} && log_info "Run all tests"
    {{ go }} test $({{ go }} list ./... | grep -Ev 'examples|cmd') -coverprofile=coverage.txt
    @. {{ logger }} && log_info "Total Coverage"
    {{ go }} tool cover -func=coverage.txt | grep total | awk '{print $3}'

# Clean go tests cache and run all tests
test-force:
    @. {{ logger }} && log_info "Clean go tests cache and run all tests"
    {{ go }} clean -testcache
    just test

# Run all tests and generate coverage report.
test-coverage:
    @. {{ logger }} && log_info "Run all tests and generate coverage report"
    {{ go }} test -count=1 -timeout 30s $({{ go }} list ./... | grep -Ev 'examples|cmd') -covermode=atomic -coverprofile=coverage.txt

# Run all tests with race detector
test-race:
    @. {{ logger }} && log_info "Running tests with race detector"
    {{ go }} test -race $({{ go }} list ./... | grep -Ev 'examples|cmd')

# === Sub-module Recipes ===

# Run tests across all sub-modules
test-all: test _test-sub

_test-sub:
    cd cobra && {{ go }} test ./...
    cd urfave && {{ go }} test ./...
    cd kong && {{ go }} test ./...

# Clean test cache and run all tests across all sub-modules
test-force-all: test-force _test-force-sub

_test-force-sub:
    cd cobra && {{ go }} clean -testcache && {{ go }} test ./...
    cd urfave && {{ go }} clean -testcache && {{ go }} test ./...
    cd kong && {{ go }} clean -testcache && {{ go }} test ./...

# Run tests with race detector across all sub-modules
test-race-all: test-race _test-race-sub

_test-race-sub:
    cd cobra && {{ go }} test -race ./...
    cd urfave && {{ go }} test -race ./...
    cd kong && {{ go }} test -race ./...

# Run lint across all sub-modules
lint-all: lint _lint-sub

_lint-sub:
    cd cobra && golangci-lint run ./...
    cd urfave && golangci-lint run ./...
    cd kong && golangci-lint run ./...

# Run fmt across all sub-modules
fmt-all: fmt
    cd cobra && {{ go }} fmt ./...
    cd urfave && {{ go }} fmt ./...
    cd kong && {{ go }} fmt ./...

# Run tidy across all sub-modules
tidy-all: tidy
    cd cobra && {{ go }} mod tidy
    cd urfave && {{ go }} mod tidy
    cd kong && {{ go }} mod tidy

# Run fmt + modernize + lint + reportcard in a single directory
_check-single dir:
    @. {{ logger }} && log_info "Checking {{ dir }}/"
    cd {{ dir }} && {{ go }} fmt ./...
    cd {{ dir }} && modernize --fix ./...
    cd {{ dir }} && golangci-lint run ./...
    cd {{ dir }} && goreportcard-cli -v

# Run full check across root + all sub-modules
check-all: check
    just _check-single cobra
    just _check-single urfave
    just _check-single kong

# === Demo Screenshots ===

freeze := "freeze"

# Capture a single style demo screenshot (dark mode only)
_capture-demo style:
    #!/usr/bin/env bash
    set -eu
    mkdir -p assets/demos
    out="assets/demos/demo-{{ style }}.png"
    HERALD_FORCE_DARK=1 {{ go }} run ./examples/demos/{{ style }}/ \
        | {{ freeze }} --output "${out}" \
            --theme "Catppuccin Mocha" --padding 20 --window \
            --width 700 \
            --shadow.blur 20 --shadow.x 0 --shadow.y 10
    # Trim trailing empty space so all screenshots share a consistent height
    magick "${out}" -trim +repage \
        -bordercolor '#1e1e2e' -border 20 \
        -gravity North \( +clone -background black -shadow 60x20+0+10 \) \
        +swap -background none -layers merge +repage \
        "${out}"

# Generate all demo screenshots and the hero collage
demo-screenshot:
    just _capture-demo compact
    just _capture-demo rich
    just _capture-demo grouped
    just _capture-demo markdown
    just _demo-hero

# Compose the 1x4 hero strip from the four style screenshots
_demo-hero:
    #!/usr/bin/env bash
    set -eu
    cd assets/demos
    gap=12
    bg='#1e1e2e'
    crop_h=500
    # Crop each to a uniform height from the top, showing key sections.
    for f in demo-compact.png demo-rich.png demo-grouped.png demo-markdown.png; do
        w=$(magick identify -format '%w' "$f")
        magick "$f" -gravity North -crop "${w}x${crop_h}+0+0" +repage "hero-${f}"
    done
    # Join horizontally with spacing.
    magick \( hero-demo-compact.png -bordercolor "$bg" -border ${gap}x${gap} \) \
           \( hero-demo-rich.png -bordercolor "$bg" -border ${gap}x${gap} \) \
           \( hero-demo-grouped.png -bordercolor "$bg" -border ${gap}x${gap} \) \
           \( hero-demo-markdown.png -bordercolor "$bg" -border ${gap}x${gap} \) \
           +append \
           -bordercolor "$bg" -border ${gap}x${gap} \
           \( +clone -background black -shadow 60x20+0+10 \) \
           +swap -background none -layers merge +repage \
           demo-hero.png
    rm -f hero-demo-*.png

# === Utilities ===

# Update dependencies
deps-update:
    @. {{ logger }} && log_info "Running go update deps"
    {{ go }} get -u ./...
    {{ go }} mod tidy
