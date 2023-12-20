# k8s-monitor-dns

This tool monitors the Kube-DNS / CoreDNS at the pod level to detect DNS issues on each node by checking internal and external DNS records. It exposes various metrics via Prometheus, including the status and response times of DNS services.

## Metrics

- `dns_status`: DNS up/down status. 1 for up, 0 for down. Labeled by 'type' which can be 'internal' or 'external'.
- `dns_response_time_seconds`: Histogram of response times for DNS checks in seconds, labeled by 'type'.

## Install

To install `k8s-monitor-dns` in your Kubernetes cluster, follow these steps:

```bash
git clone https://github.com/mattmattox/k8s-monitor-dns.git
cd k8s-monitor-dns
kubectl apply -f deploy.yaml
```

## Upgrade

To upgrade `k8s-monitor-dns` to the latest version:

```bash
git clone https://github.com/mattmattox/k8s-monitor-dns.git
cd k8s-monitor-dns
kubectl apply -f deploy.yaml
```

## Default Settings

- `InternalHost`: `kube-dns.kube-system.svc.cluster.local`
  - This record is used to check if internal DNS is working.
- `InternalIP`: `10.43.0.10`
  - This IP is used to verify the correct IP is being returned.
- `ExternalHost`: `a.root-servers.net`
  - This record is used to check if external DNS is working. Note: This should be set to a known good record that doesn't change and always returns the same IP.
- `ExternalIP`: `198.41.0.4`
  - This IP is used to verify the correct IP is being returned.
- `Timeout`: `10s`
  - This is the timeout value for DNS lookups.
- `Delay`: `5s`
  - This is the delay between checks.
