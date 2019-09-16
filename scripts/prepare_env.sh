#!/bin/bash
set -e

ssh -qT piaqua << EOF
  sudo sh -c '\
    useradd -m -d /opt/aqua -s /usr/sbin/nologin -U aqua && \
    usermod -a -G gpio aqua \
  '
EOF
