// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre;

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
public class ProjectApplication {

    public static final String SAMPLE_JSON =
        "{ \"recorded\" : 1756655999, \"location\" : \"den\", \"sensor\" : \"bmp280\", \"measurement\" : \"temperature\", \"units\" : \"C\", \"value\" : 22.3 }";

    public static String jsonToCsv(String jsonString) throws Exception {
        // Parse JSON using java.util and org.json
        org.json.JSONObject jsonObject = new org.json.JSONObject(jsonString);

        StringBuilder csvBuilder = new StringBuilder();
        StringBuilder headerBuilder = new StringBuilder();
        StringBuilder valueBuilder = new StringBuilder();

        for (String key : jsonObject.keySet()) {
            headerBuilder.append(key).append(",");
            valueBuilder.append(jsonObject.get(key).toString()).append(",");
        }

        // Remove trailing commas
        if (headerBuilder.length() > 0) {
            headerBuilder.setLength(headerBuilder.length() - 1);
        }
        if (valueBuilder.length() > 0) {
            valueBuilder.setLength(valueBuilder.length() - 1);
        }

        csvBuilder.append(headerBuilder).append("\n").append(valueBuilder);
        return csvBuilder.toString();
    }

    public static void main(String[] args) {
        SpringApplication.run(ProjectApplication.class, args);
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
        // Simply return the received JSON object
        return data;
    }

    @GetMapping("/report")
    public ResponseEntity<Resource> getReport() throws Exception {
        String csvData = ProjectApplication.jsonToCsv(
            ProjectApplication.SAMPLE_JSON
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
        // Here you can add logic to handle purge
        return Map.of("message", "Purge operation completed");
    }
}
