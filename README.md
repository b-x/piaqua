# PiAqua
Raspberry Pi aquarium controller

Under development...

## Installing

Prepare service
```sh
useradd -m -d /opt/aqua -s /usr/sbin/nologin -U aqua
usermod -a -G gpio aqua
cp aqua.service /etc/systemd/system
systemctl enable aqua.service
systemctl start aqua
```
View logs
```sh
journalctl -u aqua
```