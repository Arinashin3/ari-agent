# Storage Exporter 
Spectrum & Unisphere API를 통해 데이터를 수집하여 백엔드(ex. Otel Collector, Prometheus, Loki 등)로 성능정보를 전송한다.


## Spectrum Exporter
### 메트릭 정보

| Provider    | Default Enabled | Desc                                      |
|-------------|-----------------|-------------------------------------------|
| system      | true            | lssystem 커맨드와 동일                          |
| event       | true            | lseventlog 커맨드와 동일                        |
| performance | true            | lssystemstats 커맨드와 동일 (1m마다 최근 5s 데이터 수집) |