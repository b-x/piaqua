#!/bin/bash
set -e

scp aqua piaqua:/tmp

ssh -qT piaqua << EOF
  sudo sh -c '\
    chown aqua:aqua /tmp/aqua && \
    chmod 700 /tmp/aqua && \
    mv /tmp/aqua /opt/aqua/ && \
    systemctl restart aqua && \
    systemctl status aqua \
  '
EOF
