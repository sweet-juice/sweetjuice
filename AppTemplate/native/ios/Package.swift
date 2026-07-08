// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "sweetjuice",
    platforms: [
        .iOS(.v17),
        .macOS(.v14),
    ],
    products: [
        // An xtool project should contain exactly one library product,
        // representing the main app.
        .library(
            name: "sweetjuice",
            targets: ["sweetjuice"]
        ),
    ],
    targets: [
        .target(
            name: "sweetjuice",
            dependencies: ["Sweetjuice"]
        ),
        .binaryTarget(
            name: "Sweetjuice",
            path: "sweetjuice.xcframework"
        )
    ]
)
