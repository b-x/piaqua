# PiAqua
Raspberry Pi aquarium controller

Under development...

## Installing

Download latest Raspbian Lite from:<br>
https://www.raspberrypi.org/downloads/raspbian/

Install image to the SD card (4GB recommended)
```sh
# on linux run `lsblk` to find device file eg. `/dev/sdX`
$ umount /dev/sdX*
$ sudo dd bs=4M if=xxx.img of=/dev/sdX status=progress conv=fsync
```

### Setup WiFi
```sh
# mount sd-card (reinsert it)
$ cd /boot
$ vim wpa_supplicant.conf
```
```
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
country=«your_ISO-3166-1_two-letter_country_code»

network={
    ssid="«your_SSID»"
    psk="«your_PSK»"
    key_mgmt=WPA-PSK
}
```
```sh
# enable ssh access
$ touch ssh
```

### Configure system
```sh
# insert sd-card to raspberry pi and boot
$ ssh pi@«ip-addr»
# password: raspberry

$ sudo raspi-config
# - expand filesystem
# - change hostname
# - change password
# - set timezone
# - disable serial interface (if not needed)
# - reboot

$ sudo rpi-update
$ sudo apt update
$ sudo apt dist-upgrade

# disable bluetooth (if not needed):
$ sudo vi /boot/config.txt
# append:
dtoverlay=disable-bt
# and disable splash
disable_splash=1

# disable unused services
$ sudo systemctl disable hciuart
$ sudo systemctl disable avahi-daemon
$ sudo systemctl disable triggerhappy

# disable swap
$ sudo systemctl disable dphys-swapfile

# reduce cron tasks
$ sudo apt purge man-db

# fix wifi disconnection issue
$ sudo vi /etc/rc.local
# and append before exit:
iw wlan0 set power_save off

$ sudo reboot
```

### Enable internet access
* create dns name, eg:<br>
https://www.duckdns.org
* add cron entry to update IP (follow the site instructions)
* create NAT rules in your router:

service | inner port | outer port | protocol
------- | ---------- | ---------- | --------
web | 8080 | 80 | tcp
ssh (optional) | 22 | 2222 | tcp


### Prepare service
```sh
$ useradd -m -d /opt/aqua -s /usr/sbin/nologin -U aqua
$ usermod -a -G gpio aqua
$ cp aqua.service /etc/systemd/system
$ systemctl enable aqua.service
$ systemctl start aqua
```
View logs
```sh
$ journalctl -u aqua
```