// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre;

import com.ggeoffre.data.*;
import java.util.*;
import org.json.*;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
public class SensorApplication {

    public static SensorDataAccess getSensorDataAccess() {
        String dataAccessType = System.getenv("DATA_ACCESS");
        if (dataAccessType == null || dataAccessType.isEmpty()) {
            dataAccessType = "redis";
        }
        switch (dataAccessType.toLowerCase()) {
            case "redis":
                return new RedisData();
            case "mongo":
                return new MongoData();
            case "cassandra":
                return new CassandraData();
            case "mysql":
                return new MySQLData();
            case "postgres":
                return new PostgresData();
            default:
                throw new IllegalArgumentException(
                    "Unsupported DATA_ACCESS type: " + dataAccessType
                );
        }
    }

    public static void main(String[] args) {
        SpringApplication.run(SensorApplication.class, args);
    }
}

@RestController
@RequestMapping("/")
class SensorDataApiController {

    @GetMapping
    public Map<String, String> sayHello() {
        return Map.of("message", "Spring API Server is running!");
    }

    @PostMapping("/echo")
    public Map<String, Object> echo(@RequestBody Map<String, Object> data) {
        return data;
    }

    @PostMapping("/log")
    public Map<String, Object> log(@RequestBody Map<String, Object> data) {
        SensorDataAccess sensorData = SensorApplication.getSensorDataAccess();
        sensorData.logSensorData(new JSONObject(data).toString());
        return data;
    }

    @GetMapping("/report")
    public ResponseEntity<?> getReport() throws Exception {
        SensorDataAccess sensorData = SensorApplication.getSensorDataAccess();
        List<String> sensorDataList = sensorData.fetchSensorData();

        if (sensorDataList.isEmpty()) {
            return ResponseEntity.ok(
                Map.of("message", "No sensor data stored")
            );
        }

        String csvData = SensorDataJsonHelper.jsonArrayToCsv(
            sensorDataList.toArray(new String[0])
        );
        ByteArrayResource resource = new ByteArrayResource(csvData.getBytes());
        HttpHeaders headers = new HttpHeaders();
        headers.add(
            HttpHeaders.CONTENT_DISPOSITION,
            "attachment; filename=report.csv"
        );
        headers.add(HttpHeaders.CONTENT_TYPE, "text/csv");
        return ResponseEntity.ok()
            .headers(headers)
            .contentLength(csvData.length())
            .body(resource);
    }

    @RequestMapping(
        value = "/purge",
        method = { RequestMethod.GET, RequestMethod.POST }
    )
    public Map<String, String> purge() {
        SensorDataAccess sensorData = SensorApplication.getSensorDataAccess();
        sensorData.purgeSensorData();
        return Map.of("message", "Purge operation completed");
    }
}
