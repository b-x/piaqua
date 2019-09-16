#!/bin/bash
set -e

scp -r public piaqua:/tmp

ssh -qT piaqua << EOF
  sudo sh -c '\
    chown -R aqua:aqua /tmp/public && \
    chmod 444 /tmp/public/* && \
    rm -rf /opt/aqua/public && \
    mv /tmp/public /opt/aqua/ \
  '
EOF
