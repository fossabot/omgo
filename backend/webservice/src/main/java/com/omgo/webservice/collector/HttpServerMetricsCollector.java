package com.omgo.webservice.collector;

import io.prometheus.client.Collector;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.dropwizard.MetricsService;

import java.util.ArrayList;
import java.util.List;

public class HttpServerMetricsCollector extends Collector {
    private Vertx vertx;
    private MetricsService metricsService;
    private HttpServer server;

    HttpServerMetricsCollector(Vertx vertx, MetricsService metricsService, HttpServer server) {
        this.vertx = vertx;
        this.metricsService = metricsService;
        this.server = server;
    }

    private void addHttpServerMetrics(List<MetricFamilySamples> sampleFamilies) {
        JsonObject metrics = metricsService.getMetricsSnapshot(server);
        if (metrics.containsKey("get-requests")) {
            JsonObject dataJson = metrics.getJsonObject("get-requests");
            sampleFamilies.add(createHistogramMFS("get-request", "rate of http get method occurrence", dataJson));
        }

    }

    @Override
    public List<MetricFamilySamples> collect() {
        List<MetricFamilySamples> mfs = new ArrayList<MetricFamilySamples>();
        addHttpServerMetrics(mfs);
        return null;
    }

    private MetricFamilySamples createHistogramMFS(String name, String help, JsonObject data) {
        List<MetricFamilySamples.Sample> samples = new ArrayList<>();

        MetricFamilySamples.Sample sample = new MetricFamilySamples.Sample("count", null, null, data.getLong("count"));
        data.getLong("count");
        return new MetricFamilySamples(name, Type.HISTOGRAM, help, samples);
    }
}
