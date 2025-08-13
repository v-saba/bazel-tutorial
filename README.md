# Bazel + Go Tutorial Project

This project demonstrates a modern Bazel setup with Go using Bzlmod (MODULE.bazel) instead of the legacy WORKSPACE system.

## Project Structure

```
bazel-tutorial/
├── MODULE.bazel          # Bazel module configuration (replaces WORKSPACE)
├── MODULE.bazel.lock     # Dependency lockfile (commit to git)
├── BUILD                 # Root BUILD file with Gazelle target
├── go.mod               # Go module file (external deps only)
├── go.sum               # Go checksums (external deps only)
├── proto/               # Protocol buffer definitions
│   ├── BUILD
│   └── v1/
│       └── telemetry_server.proto
├── server/              # Server binary
│   ├── BUILD
│   └── main.go
├── client/              # Client binary (if exists)
│   ├── BUILD
│   └── main.go
└── common/              # Shared internal library
    ├── BUILD
    └── common.go
```

## Initial Project Setup

### 1. Initialize the project
```bash
# Create go.mod for external dependencies only
go mod init github.com/your-org/your-project

# Create MODULE.bazel with basic configuration
cat > MODULE.bazel << 'EOF'
module(name = "your_project")

bazel_dep(name = "rules_proto", version = "7.0.2")
bazel_dep(name = "rules_go", version = "0.46.0")
bazel_dep(name = "gazelle", version = "0.35.0")

# Go SDK
go = use_extension("@rules_go//go:extensions.bzl", "go_sdk")
go.download(version = "1.21.0")

# Go dependencies from go.mod
go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//:go.mod")

# Import external dependencies
use_repo(go_deps)
EOF

# Create root BUILD file with Gazelle
cat > BUILD << 'EOF'
load("@gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")
EOF
```

### 2. Generate initial lockfile
```bash
bazel mod deps
```

## Working with Dependencies

### External Dependencies (3rd party packages)

#### Adding a new external dependency:
```bash
# 1. Add to go.mod
go get github.com/some/package@v1.2.3

# 2. Update MODULE.bazel use_repo section
# Add the repository name to use_repo:
use_repo(
    go_deps,
    "com_github_some_package",  # Add this line
)

# 3. Create go.sum entry
go mod download

# 4. Reference in BUILD files
deps = [
    "@com_github_some_package//:go_default_library",
]
```

#### Repository naming convention:
- `github.com/google/uuid` → `com_github_google_uuid`
- `golang.org/x/net` → `org_golang_x_net`
- `google.golang.org/grpc` → `org_golang_google_grpc`

### Internal Dependencies

#### Never add internal packages to go.mod - handle via Bazel only!

```bash
# Create new internal library
mkdir -p internal/auth
cat > internal/auth/BUILD << 'EOF'
load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "auth",
    srcs = ["auth.go"],
    importpath = "github.com/your-org/your-project/internal/auth",
    visibility = ["//visibility:public"],
)
EOF

# Reference in other BUILD files
deps = [
    "//internal/auth",  # Internal dependency
]
```

## Adding New Libraries/Binaries

### Go Library
```bash
mkdir -p mylib
cat > mylib/BUILD << 'EOF'
load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "mylib",
    srcs = ["mylib.go"],
    importpath = "github.com/your-org/your-project/mylib",
    visibility = ["//visibility:public"],
    deps = [
        "//common",  # Internal deps
        "@com_github_google_uuid//:go_default_library",  # External deps
    ],
)

go_test(
    name = "mylib_test",
    srcs = ["mylib_test.go"],
    embed = [":mylib"],
)
EOF
```

### Go Binary
```bash
mkdir -p mycmd
cat > mycmd/BUILD << 'EOF'
load("@rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "mycmd",
    srcs = ["main.go"],
    deps = [
        "//mylib",
        "//common",
    ],
)
EOF
```

### Protocol Buffers
```bash
mkdir -p proto/v2
cat > proto/v2/my_service.proto << 'EOF'
syntax = "proto3";
package myservice.v2;
option go_package = "github.com/your-org/your-project/proto/gen/go/myservice/v2";
// ... proto definition
EOF

cat > proto/BUILD << 'EOF'
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "myservice_proto",
    srcs = ["v2/my_service.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "myservice_go_proto",
    importpath = "github.com/your-org/your-project/proto/gen/go/myservice/v2",
    proto = ":myservice_proto",
    visibility = ["//visibility:public"],
)
EOF
```

## Common Commands

### Building
```bash
# Build everything
bazel build //...

# Build specific target
bazel build //server:server
bazel build //mylib:mylib

# Build and run
bazel run //server:server
```

### Testing
```bash
# Test everything
bazel test //...

# Test specific package
bazel test //mylib:mylib_test
```

### Dependency Management
```bash
# Update dependency lockfile
bazel mod deps

# Clean build cache
bazel clean

# Query dependencies
bazel query "deps(//server:server)"
bazel query "//..."
```

### Gazelle (Auto-generate BUILD files)
```bash
# Update BUILD files (if using Gazelle for auto-generation)
bazel run //:gazelle

# Update go dependencies
bazel run //:gazelle -- update-repos -from_file=go.mod
```

## Development Workflow

### 1. Adding a new feature with external dependency
```bash
# Add external dependency
go get github.com/new/package@v1.0.0

# Update MODULE.bazel
# Add to use_repo section: "com_github_new_package"

# Create your code
mkdir -p myfeature
# Write myfeature/main.go and myfeature/BUILD

# Build and test
bazel build //myfeature:myfeature
bazel test //myfeature:myfeature_test
```

### 2. Updating dependencies
```bash
# Update go.mod
go get -u github.com/some/package

# Update go.sum
go mod download

# Rebuild
bazel build //...
```

### 3. Refactoring internal code
```bash
# Move/rename packages - update BUILD files
# Update import paths in BUILD deps
# No changes needed in go.mod for internal packages
```

## File Management Rules

### Commit to Git:
- ✅ `MODULE.bazel`
- ✅ `MODULE.bazel.lock`
- ✅ `go.mod`
- ✅ `go.sum`
- ✅ All `BUILD` files
- ✅ Source code (`.go`, `.proto`)

### Ignore in Git:
- ❌ `bazel-*` symlinks/directories
- ❌ Build artifacts

### Key Principles:
1. **go.mod is for external dependencies only** - never add internal packages
2. **BUILD files handle all internal dependencies** via `//package:target` syntax
3. **MODULE.bazel.lock ensures reproducible builds** - always commit it
4. **use_repo must list all external Go dependencies** you reference in BUILD files
5. **Repository names follow conversion rules** (dots to underscores, domain reversal)

## Troubleshooting

### Common Issues:
1. **"No repository visible as '@com_github_...'"**
   - Add the repository to `use_repo()` in MODULE.bazel

2. **"Cannot find module providing package"**
   - Don't add internal packages to go.mod
   - Use BUILD file dependencies for internal code

3. **Gazelle segfault**
   - Usually due to invalid go.mod or missing go.sum
   - Ensure go.mod only has external dependencies

4. **Build cache issues**
   - Run `bazel clean` to reset build state
