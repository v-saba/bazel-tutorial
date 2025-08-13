load("@rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "server",
    srcs = ["main.go"],
    deps = [
        "//proto:go_default_library",
        "//common:go_default_library"
    ],
)
