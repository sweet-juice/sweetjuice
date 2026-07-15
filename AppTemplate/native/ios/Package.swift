// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "GenericApp",
    platforms: [
        .iOS(.v17),
        .macOS(.v14),
    ],
    products: [
        .library(
            name: "App",
            targets: ["GenericApp"]
        ),
    ],
    targets: [
        .target(
            name: "GenericApp",
            dependencies: ["Sweetjuice"]
        ),
        .binaryTarget(
            name: "Sweetjuice",
            path: "Sweetjuice.xcframework"
        )
    ]
)
