// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre;

import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.json.JSONObject;

@Path("/")
public class RootResource {

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

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public String root() {
        return "{\"message\": \"Quarkus API Server is running!\"}";
    }

    @POST
    @Path("/echo")
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public String echo(String json) {
        return json;
    }

    @POST
    @Path("/log")
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public String log(String logMessage) {
        return logMessage;
    }

    @GET
    @Path("/report")
    @Produces(MediaType.TEXT_PLAIN)
    public jakarta.ws.rs.core.Response report() {
        try {
            String csvData = jsonToCsv(SAMPLE_JSON);
            return jakarta.ws.rs.core.Response.ok(csvData)
                .header(
                    "Content-Disposition",
                    "attachment; filename=\"report.csv\""
                )
                .header("Content-Type", "text/csv")
                .build();
        } catch (Exception e) {
            return jakarta.ws.rs.core.Response.status(
                jakarta.ws.rs.core.Response.Status.INTERNAL_SERVER_ERROR
            )
                .entity("Error generating CSV: " + e.getMessage())
                .build();
        }
    }

    @Path("/purge")
    @Consumes(MediaType.TEXT_PLAIN)
    @Produces(MediaType.APPLICATION_JSON)
    @GET
    @POST
    public String purge() {
        return "{\"message\": \"Purge operation executed\"}";
    }
}
