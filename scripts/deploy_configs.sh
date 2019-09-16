#!/bin/bash
set -e

scp -r configs piaqua:/tmp

ssh -qT piaqua << EOF
  sudo sh -c '\
    chown -R aqua:aqua /tmp/configs && \
    chmod 444 /tmp/configs/* && \
    mv /tmp/configs /opt/aqua/ \
  '
EOF
