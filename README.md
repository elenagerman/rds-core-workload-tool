# Testcmd

## Overview

Simple test application which allows to run icmp/tcp/udp/sctp connectivity test using unicast/multicast/broadcast

In order to build testcmd application run:

`make testcmd-bin`

## Supported Protocols

#### IPV4

* **icmp**
* **sctp**
* **tcp**
* **udp**
    * **multicast**
    * **broadcast**

#### IPV6

* **icmp**
* **sctp**
* **tcp**
* **udp**
    * **multicast**

## Flags

* **listen** - insert this flag in order to run server
* **interface** - insert this flag to specify the interface you want to use(Examples: ens33/eth0/net1)
* **multicast** - insert this flag in order to run a udp **multicast** server
* **broadcast** - insert this flag in order to run a udp **broadcast** server
* **protocol** -  protocol name (Options: tcp/udp/icmp/sctp)
* **mtu** - MTU size. Any integer number in range 50-9000 (deafult 1450)
* **server** - destination IPv4/IPv6 address
* **port** - port number. Any integer number in range 1-65534 (default 80)
* **negative** - insert this flag if **no** connectivity is expected
* **packages** - packages number. Any integer number in range 1-65534 (default 5)
* **timeoutTCP** - session timeout. Any integer number in range 1-65534 (default 2)
* **timeoutUDP** - session timeout. Any integer number in range 1-65534 (default 5)
