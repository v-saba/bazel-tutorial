load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_go//go:def.bzl", "go_library")
load("@rules_go//proto:def.bzl", "go_proto_library")

# Proto library for telemetry server
proto_library(
    name = "telemetry_server_proto",
    srcs = ["v1/telemetry_server.proto"],
    visibility = ["//visibility:public"],
)

# Go proto library with gRPC
go_proto_library(
    name = "telemetry_server_go_proto",
    compilers = ["@rules_go//proto:go_grpc"],
    importpath = "github.com/v-saba/bazel-tutorial/proto/gen/go/telemetry_server/v1",
    proto = ":telemetry_server_proto",
    visibility = ["//visibility:public"],
)

# Go library for easier imports
go_library(
    name = "go_default_library",
    embed = [":telemetry_server_go_proto"],
    importpath = "github.com/v-saba/bazel-tutorial/proto/gen/go/telemetry_server/v1",
    visibility = ["//visibility:public"],
)
