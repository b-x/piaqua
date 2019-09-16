#!/bin/bash
set -e

scp scripts/aqua.service piaqua:/tmp

ssh -qT piaqua << EOF
  sudo sh -c '\
    chown root:root /tmp/aqua.service && \
    mv /tmp/aqua.service /etc/systemd/system && \
    systemctl enable aqua.service \
  '
EOF
