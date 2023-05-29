
wget -c https://github.com/prometheus/prometheus/releases/download/v2.28.1/prometheus-2.28.1.linux-amd64.tar.gz
tar -vxzf prometheus-2.28.1.linux-amd64.tar.gz
mv prometheus-2.28.1.linux-amd64 /usr/local/prometheus

systemctl daemon-reload
systemctl start prometheus.service