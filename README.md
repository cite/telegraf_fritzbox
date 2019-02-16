# telegraf\_fritzbox

## Overview

This is a **very** quick and dirty attempt at enabling
[telegraf](https://github.com/influxdata/telegraf/) to collect basic data from
the popular FRITZ!Box routers made by manufacturer AVM. These devices export
basic statistcis using UPnP (i.e. SOAP on port 49000) when configured to do so
("Heimnetz > Netzwerk > Netzwerkeinstellungen > Statusinformationen über UPnP
übertragen").

The code borrows heavily from Nils Decker's excellent
[fritzbox\_exporter](https://github.com/ndecker/fritzbox_exporter).

## Build

1. Copy code to `$TELEGRAF_SOURCE/plugins/all/fritzbox/fritzbox.go`, creating
   directories as needed.
2. Add `"github.com/influxdata/telegraf/plugins/inputs/fritzbox"` to imports in
   `plugins/inputs/all/all.go`.
3. Build as usual.

The code doesn't handle data types, and I've never tested it with either
multiple FRITZ!Box devices or FRITZ!OS version other than 7.02.

## Status

Works for me. Contributions welcome.
