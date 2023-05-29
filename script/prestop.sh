#!/usr/bin/env bash
echo $CLOUD_POD_IP
echo $ETCDCTL_ENDPOINTS
echo "/psm/$CLOUD_POD_IP:9999"
etcdctl del "/psm/$CLOUD_POD_IP:9999"