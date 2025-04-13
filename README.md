# OpenEVT - Envertec EVT800 Client

Take control of your solar energy monitoring with OpenEVT!

OpenEVT is tool that empowers you to directly communicate with your Envertec
solar microinverter, giving you complete ownership of your solar performance
data - without being restricted to the Envertec APIs. OpenEVT connects directly
to your microinverter on your LAN.

- Improved Stability: Continue monitoring your PV performance even during an
  Envertec Cloud outage.
- Data Privacy: Be in charge of your data. Monitor your PV performance within
  your own private network.
- Integration Ready: Integrate PV monitoring into Home Assistant, Grafana and
  other monitoring tools and home automation systems.

OpenEVT is known to work with the following Envertec inverters:

- `EVT800B`

## Install

To install OpenEVT (Go 1.21+):

```shell
$ go install https://github.com/brandon1024/OpenEVT/cmd/openevt@latest
```

## Usage

To connect to your microinverter, you need:

- the LAN address and port number (e.g. `192.0.2.1:14889`),
- the inverter serial number (e.g. `31583078`)

Before connecting to your inverter, you must set up the inverter following the
instructions provided by the manufacturer. The inverter enters a low-power
standby mode when there's no sunlight, so OpenEVT won't be able to connect
during the night.

```shell
$ openevt --addr 192.168.2.54:14889 --serial-number 31583078
```

For some usage info:

```shell
$ openevt --help
  -addr string
        address and port of the microinverter (e.g. 192.0.2.1:14889)
  -reconnect-interval duration
        interval between connection attempts (e.g. 1m) (default 1m0s)
  -serial-number string
        serial number of your microinverter (e.g. 31583078)
  -web.disable-exporter-metrics
        exclude metrics about the exporter itself (go_*)
  -web.listen-address string
        address on which to expose metrics (default ":9090")
  -web.telemetry-path string
        path under which to expose metrics (default "/metrics")
```

### Finding your Inverter on the LAN

To find the address and port of your inverter, connect to the wireless access
point of your inverter. The SSID is the serial number of your inverter, like
`31583078`. Once connected, login using your credentials.

Under the `System` tab, the LAN IP address can be found in the `STA Mode` `IP
Address` field.

Under the `Other Settings` tab, you can find the port number in the `Port ID`
field.

Alternatively, you can scan your network with a tool like nmap:

```shell
nmap -p 14889 192.168.2.0/24
```

## Building

To build OpenEVT (Go 1.21+):

```shell
make
```

### Contributing

Help us support more inverter models! If your inverter also supports a local
mode connection, send us your packet captures to help us expand support for more
inverters.

If OpenEVT doesn't work for your particular inverter model, please [create an
issue](https://github.com/brandon1024/OpenEVT/issues) and we'll do our best to
support you.

## Technical Info

OpenEVT was developed by reverse engineering the communication between the
EnverView app and the microinverter's local mode port. There's a lot we don't
understand yet, but here's what we've found so far.

### Poll Message Format

When we first connect to the inverter, we issue a poll message to the inverter
to request it's current state. The microinverter will immediately respond with
an _inverter state_ message. The message is 32 bytes long and has the following
format:

```
RAW:    3638 3030 3130 3638 3130 3737 3332 3332 3332 3332 3030 3030 3030 3030 3966 3136
ASCII:   68   00   10   68   10   77   32   32   32   32   00   00   00   00   9f   16
        |----------------------|      |-----------------| |-----------------|      |--|
                 FIXED                |   INVERTER SN   | |     PADDING     |      DONE
```

This message looks quite similar to the _acknowledge_ message, the key
differences being the `77` and `9f` words.

### Inverter State Message Format

Periodically, the inverter will push a message to the client that contains
performance metrics of both inverter modules. The message is 86 bytes long and
has the following format:

```
RAW: 6800 5668 1051 3232 3232 7001 7900 0000 0000 0000 3232 3232 7079 2f47 00e6 0003 49dc 22b3 3a7e 31f8 0200 0000 0000 0000 0000 0000 3232 3233 7079 302b 00e4 0002 c7f6 2319 3a7e 31f8 0200 0000 0000 0000 0000 0000 9b16
                    |-------| |-------| |------------| |-------| |--| |--| |--| |-------| |--| |--| |--| |---------------------------| |-------| |--| |--| |--| |-------| |--| |--| |--| |---------------------------|   ||
                      EVTID      FW?                      MID     FW   #1   #2     #3      #4   #5   #6                                   MID     FW   #1   #2     #3      #4   #5   #6
```

Metrics for the first module are decoded as follows:

- #1: [26,27] DC Input Voltage:    23.64V     `(* 64 / 32768)`
- #2: [28,29] AC Output Power:     3.59W      `(* 512 / 32768)`
- #3: [30-33] Total Energy:        26.31kWh   `(* 4 / 32768)`
- #4: [34,35] Temperature:         29.40C     `(* (256 / 32768) - 40)`
- #5: [36,37] AC Output Voltage:   233.97V    `(* 512 / 32768)`
- #6: [38,39] AC Output Frequency: 49.97Hz    `(* 128 / 32768)`

Metrics for the second module are decoded as follows:

- #1: [58,59] DC Input Voltage:    24.08V     `(* 64 / 32768)`
- #2: [60,61] AC Output Power:     3.56W      `(* 512 / 32768)`
- #3: [62-65] Total Energy:        22.25kWh   `(* 4 / 32768)`
- #4: [66,67] Temperature:         30.20C     `(* (256 / 32768) - 40)`
- #5: [68,69] AC Output Voltage:   233.97V    `(* 512 / 32768)`
- #6: [70,71] AC Output Frequency: 49.97Hz    `(* 128 / 32768)`

### Acknowledge Message Format

Once we receive an _inverter state_ message from the inverter, we need to
acknowledge it with an acknowledge message. The message format is quite similar
to the _poll_ message (the `50` and `78` words are relevant):

```
HEX:    3638 3030 3130 3638 3130 3530 3332 3332 3332 3332 3030 3030 3030 3030 3738 3136
ASCII:   68   00   10   68   10   50   32   32   32   32   00   00   00   00   78   16
        |----------------------|      |-----------------| |-----------------|      |--|
                                      |   INVERTER SN   | |     PADDING     |
```

It was found that the client needs to acknowledge messages otherwise the
inverter hangs up the connection and disconnects the client.

## License

MIT License. Copyright (c) 2022 Brandon Richardson.
