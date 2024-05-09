# nosana-exporter

For [Nosana.io](https://nosana.io/) GPU nodes.

Prometheus exporter reporting data from your Nosana.io node:

- Number of jobs running

Prometheus metrics available at (url: `/metrics`),

_NOTE: This is a work in progress. Output format is subject to change._

### Usage

```
docker run -d --restart always -p 8995:8995 bengivre/nosana-exporter \
    --podman-url=http://127.0.0.1:8080/v3.4.2/libpod/containers/json
```

You can also give other arguments or default values will be used :

```
docker run -d --restart always -p 8995:8995 bengivre/nosana-exporter \
    --podman-url=http://127.0.0.1:8080/v3.4.2/libpod/containers/json \
    --server-address=192.168.0.20 \
    --server-port=8995
```

# Grafana

![grafana img](https://github.com/bengivre/nosana-exporter/blob/main/img/grafana.png?raw=true)
