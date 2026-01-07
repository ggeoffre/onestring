// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre;

import com.ggeoffre.data.*;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import java.util.*;
import org.json.*;

@Path("/")
public class RootResource {

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
    public String log(String data) {
        SensorDataAccess sensorData = RootResource.getSensorDataAccess();
        sensorData.logSensorData(new JSONObject(data).toString());
        return data;
    }

    @GET
    @Path("/report")
    @Produces({ MediaType.TEXT_PLAIN, MediaType.APPLICATION_JSON })
    public jakarta.ws.rs.core.Response report() {
        try {
            SensorDataAccess sensorData = RootResource.getSensorDataAccess();
            List<String> sensorDataList = sensorData.fetchSensorData();
            if (sensorDataList.isEmpty()) {
                return jakarta.ws.rs.core.Response.status(
                    Response.Status.NOT_FOUND
                )
                    .entity("{\"message\": \"No sensor data stored\"}")
                    .build();
            }
            String csvData = SensorDataJsonHelper.jsonArrayToCsv(
                sensorDataList.toArray(new String[0])
            );
            return jakarta.ws.rs.core.Response.ok(csvData)
                .header(
                    "Content-Disposition",
                    "attachment; filename=\"report.csv\""
                )
                .header("Content-Type", "text/csv")
                .build();
        } catch (Exception e) {
            return jakarta.ws.rs.core.Response.status(
                Response.Status.INTERNAL_SERVER_ERROR
            )
                .entity(
                    "{\"message\": \"An error occurred while generating the report\"}"
                )
                .build();
        }
    }

    @GET
    @POST
    @Path("/purge")
    @Produces({ MediaType.APPLICATION_JSON, MediaType.TEXT_PLAIN })
    @Consumes(
        { MediaType.TEXT_PLAIN, MediaType.APPLICATION_JSON, MediaType.WILDCARD }
    )
    public String purge() {
        SensorDataAccess sensorData = RootResource.getSensorDataAccess();
        sensorData.purgeSensorData();
        return "{\"message\": \"Purge operation executed\"}";
    }
}
