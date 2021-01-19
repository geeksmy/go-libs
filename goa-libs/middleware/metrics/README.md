### 使用方法示例

1. 在 接口服务启动前开启指标导出服务

```go
import metricsMlwr "starter/middleware/metrics"

var metrics *metricsMlwr.Prometheus
if config.C.Metrics.Enabled {
    metrics = metricsMlwr.NewPrometheus("starter", nil)
    go metrics.Start(config.C.Metrics.Addr)
}

// pprof 这里启动失败会直接 panic
if config.C.Pprof.Enabled {
    go starterCmd.RunDebugPprofServer(config.C.Pprof.Addr)
}

starterCmd.RunServer(config.C, metrics)
```

2. 在 `handleHTTPServer()` 方法中使用此中间件替换 `goa` 的 `httpmdlwr.Log`方法

```go
var handler http.Handler = mux
{
    if debug {
        handler = httpmdlwr.Debug(mux, os.Stdout)(handler)
    }

    handler = mdlwr.PopulateRequestContext()(handler)
    handler = httpmdlwr.RequestID()(handler)

    if metrics != nil {
        handler = metrics.HandlerFunc(adapter)(handler)
    } else {
        handler = httpmdlwr.Log(adapter)(handler)
    }
}
```
