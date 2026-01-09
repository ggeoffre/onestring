// swift-tools-version:5.9
// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import PackageDescription

let package = Package(
    name: "AccessApp",
    platforms: [
        .macOS(.v13)  // Specify the minimum macOS version
    ],
    products: [
        // Define the executable product
        .executable(name: "AccessApp", targets: ["Run"])
    ],
    dependencies: [
        // Add any external dependencies here (if needed)
    ],
    targets: [
        // Define the library target for the Access module
        .target(
            name: "Access",
            dependencies: [],
            path: "Sources/Access"
        ),
        // Define the executable target for the Run module
        .executableTarget(
            name: "Run",
            dependencies: ["Access"],
            path: "Sources/Run"
        ),
    ]
)
