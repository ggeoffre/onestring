// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

public protocol SensorDataAccess {
    /// Logs sensor data
    func logSensorData(jsonData: String) throws

    /// Fetches sensor data
    func fetchSensorData() throws -> [String]

    /// Purges sensor data
    func purgeSensorData() throws
}
