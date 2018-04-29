# Disable Sonos WiFi
Disable sonos built in WiFi (and bridge) to prevent configuration mistakes.

Or just open the page: http://SONOS_IP:1400/wifictrl and choose "Persist Off".

### Default adapter names
* darwin = en0
* linux = eth0
* windows = Ethernet

### Build
* GOOS=windows GOARCH=amd64 go build
https://golang.org/doc/install/source#environment
