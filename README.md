# esptool

> This is a work in progress and has not yet been fully tested
> I bundle my [fluepdot](https://github.com/fluepke/fluepdot) firmware with this tool to ship a standalone update utility, which I prefer over ~~remote code execution on customer HW~~ OTA updates.

## Compilation

```bash
go get github.com/fluepke/esptool
cd ${GOPATH-$HOME/go}/src/github.com/fluepke/esptool
go build
./esptool <args>
```

## Usage

> It only works on **Linux** as of right now, because Microsoft(r) Windows(tm)(c) does not support termios. Concerning MacOS: Buy me a macbook.

Esptool offers the following subcommands:
  * version: Show version info and exit
  * info: Retrieve various information from chip
  * flashRead: Read flash contents
  * flashWrite: Write flash contents

to see the help, type `./esptool <subcommand> -h`

### Examples

Read various information and the partition table from chip, then display them in **JSON** format:
```bash
./esptool info -serial.port /dev/ttyUSB0 -json
```

<details>
  <summary>Click to see command output</summary>

  ```json
  {
    "ChipType": "ESP32D0WDQ6",
    "Revision": "1",
    "Features": [
      "240MHz",
      "WiFi",
      "Single Core",
      "VRef calibration in efuse",
      "Coding Scheme None",
      "Bluetooth"
    ],
    "MacAddress": "24:6f:28:92:ef:20",
    "Partitions": [
      {
        "name": "nvs",
        "type": "data",
        "subtype": "nvs",
        "Offset": 36864,
        "size": 16384
      },
      {
        "name": "otadata",
        "type": "data",
        "subtype": "factory",
        "Offset": 53248,
        "size": 8192
      },
      {
        "name": "phy_init",
        "type": "data",
        "subtype": "phy",
        "Offset": 61440,
        "size": 4096
      },
      {
        "name": "factory",
        "type": "app",
        "subtype": "factory",
        "Offset": 65536,
        "size": 3145728
      },
      {
        "name": "config",
        "type": "66",
        "subtype": "35",
        "Offset": 3211264,
        "size": 4096
      }
    ]
  }
  ```
</details>

Write data to flash
```bash
./esptool flashWrite -flash.file=/home/fluepke/git/fluepdot/software/firmware/flipdot-firmware.bin -flash.offset=0x10000 -serial.port=/dev/ttyUSB0 -serial.baudrate.transfer=500000 -serial.baudrate.connect=115200
```
