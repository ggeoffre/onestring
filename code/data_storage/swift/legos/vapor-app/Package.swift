// swift-tools-version:5.9
// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import PackageDescription

let package = Package(
    name: "vapor-app",
    platforms: [
        .macOS(.v10_15)
    ],
    dependencies: [
        .package(url: "https://github.com/vapor/vapor.git", from: "4.119.2"),
        .package(url: "https://github.com/apple/swift-cassandra-client", from: "0.8.0"),
        .package(url: "https://github.com/orlandos-nl/MongoKitten.git", from: "7.10.0"),
        .package(url: "https://github.com/Mordil/RediStack.git", from: "1.6.2"),
        .package(url: "https://github.com/vapor/postgres-nio.git", from: "1.29.0"),
        .package(url: "https://github.com/vapor/mysql-nio.git", from: "1.9.0"),
        .package(url: "https://github.com/apple/swift-log.git", from: "1.8.0"),
        .package(url: "https://github.com/apple/swift-nio.git", from: "2.92.1"),
        .package(url: "https://github.com/apple/swift-collections.git", from: "1.0.0"),
    ],
    targets: [
        .executableTarget(
            name: "App",
            dependencies: [
                .product(name: "Vapor", package: "vapor"),
                .product(name: "CassandraClient", package: "swift-cassandra-client"),
                .product(name: "MongoKitten", package: "MongoKitten"),
                .product(name: "RediStack", package: "RediStack"),
                .product(name: "PostgresNIO", package: "postgres-nio"),
                .product(name: "MySQLNIO", package: "mysql-nio"),
                .product(name: "Logging", package: "swift-log"),
                .product(name: "NIOCore", package: "swift-nio"),
                .product(name: "NIOPosix", package: "swift-nio"),
                .product(name: "NIOConcurrencyHelpers", package: "swift-nio"),
                .product(name: "DequeModule", package: "swift-collections"),
            ],
            path: "Sources"
        )
    ]
)
