# k8s-monitor-dns
This tools will monitor the Kube-DNS / CoreDNS at the pod level to detect DNS issues.

## Install
```
git clone https://github.com/mattmattox/k8s-monitor-dns.git
cd k8s-monitor-dns
kubectl apply -f .
```

## Upgrade
```
rm -rf k8s-monitor-dns
git clone https://github.com/mattmattox/k8s-monitor-dns.git
cd k8s-monitor-dns
kubectl apply -f .
```

## Default settings
InternalHost = kube-dns.kube-system.svc.cluster.local (This record is used to check if internal DNS is working)

InternalIP = 10.43.0.10 (This IP is used to verify the correct IP is being returned)

ExternalHost = a.root-servers.net (This record is used to check if external DNS is working. NOTE: This should be set to a known good record that doesn't change and always returns the same IP)

ExternalIP = 198.41.0.4 (This IP is used to verify the correct IP is being returned)

Timeout = 10s (This is the timeout value for DNS lookups)

Delay = 5s (This is the delay between checks)
