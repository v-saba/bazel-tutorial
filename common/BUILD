load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["common.go"],
    deps = [
        "@com_github_google_uuid//:go_default_library",
    ],
    importpath = "github.com/v-saba/bazel-tutorial/common",
    visibility = ["//visibility:public"],
)