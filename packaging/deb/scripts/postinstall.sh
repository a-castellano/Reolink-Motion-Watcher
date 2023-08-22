#!/bin/sh

mkdir -p /etc/windmaker-reolink-motion-watcher

echo "### NOT starting on installation, please execute the following statements to configure windmaker-reolink-motion-watcher to start automatically using systemd"
echo " sudo /bin/systemctl daemon-reload"
echo " sudo /bin/systemctl enable windmaker-reolink-motion-watcher"
echo "### You can start grafana-server by executing"
echo " sudo /bin/systemctl start windmaker-reolink-motion-watcher"
