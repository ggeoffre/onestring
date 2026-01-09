// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

protocol SensorDataDatabase {
    func connect() async throws
    func insert(sensorData: SensorData) async throws
    func selectAll() async throws -> [SensorData]
    func deleteAll() async throws
    func close() async throws
}
