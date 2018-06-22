# Directory Watcher - NR Infrastructure On-Host Integration
Reports contents of a specified directory to NRI Inventory

### [Download The Latest Release HERE](https://source.datanerd.us/FIT/DirWatcher/releases/latest)

### Screenshots
![alt text](https://source.datanerd.us/FIT/DirWatcher/blob/master/images/DirWatcher.jpg "Super Cool Screenshot of DirWatcher entries in Infra Inventory")

### Requirements
* NRI Agent

### NRI Configuration

* Install 3 files noted below in the 3 seperate folders
  * `dirwatcher` or `dirwatcher.exe` executable located in Releases
  * `newrelic-infra-dirwatcher-config.yml` & `newrelic-infra-dirwatcher-definition.yml` [created per the examples below](#example-config-files)
* Resart NRI Agent
* Verify data in NewRelic under NRI Inventory

### File Structure

##### Windows
* `C:\Program Files\New Relic\newrelic-infra\integrations.d\newrelic-infra-dirwatcher-config.yml`
* `C:\Program Files\New Relic\newrelic-infra\custom-integrations\newrelic-infra-dirwatcher-definition.yml`
* `C:\Program Files\New Relic\newrelic-infra\custom-integrations\newrelic-infra-dirwatcher\dirwatcher.exe`

##### Linux
* `/etc/newrelic-infra/integrations.d/newrelic-infra-dirwatcher-config.yml`
* `/var/db/newrelic-infra/custom-integrations/newrelic-infra-dirwatcher-definition.yml`
* `/var/db/newrelic-infra/custom-integrations/newrelic-infra-dirwatcher/dirwatcher`

### Example config files

#### newrelic-infra-dirwatcher-definition.yml
Interval currently set to run at 1 minute.
```
#
# New Relic Infrastructure DirWatcher Integration
#
name: com.newrelic.dirwatcher
description: DirWatcher
protocol_version: 1
commands:
  inventory:
    command:
      - .\newrelic-infra-dirwatcher\dirwatcher.exe
      - --inventory
    prefix: config/feature-flags
    interval: 3600

```
#### newrelic-infra-dirwatcher-config.yml
```
#
# New Relic Infrastructure DirWatcher Integration
#
integration_name: com.newrelic.dirwatcher
instances:
  - name: dirwatcher-inventory
    command: inventory
    arguments:
      dir_name: C:\crdata\mutex
      do_recurse: true
```

### Testing Locally

* You can run the NRI Integration Exectuable locally by manually passing in the required arguments.
* This can help when testing a new feature, or prior to running as an NRI Integration on a new server.

#### Example Command
```
/var/db/newrelic-infra/custom-integrations/newrelic-infra-dirwatcher/dirwatcher -dir_name /tmp -pretty -inventory -do_recurse true
```
#### Output
You can ignore WARN messages like these when testing locally as it is expected:
```
WARN[0000] Environment variable NRIA_CACHE_PATH is not set, using default /tmp/dirwatcher.json
WARN[0000] Cache file (/tmp/dirwatcher.json) is older than 1m0s, skipping loading from disk.
```
The command should return a JSON blob similar to the one seen here:
```json
{
  "name": "com.newrelic.dirwatcher",
  "protocol_version": "1",
  "integration_version": "0.1.0",
  "metrics": [],
  "inventory": {
    "/tmp/": {
      "fileSize": 4096,
      "isDir": true,
      "modTime": "2018-06-22T12:55:01.200689498-07:00",
      "mode": "dtrwxrwxrwx"
    },
    "/tmp/.ICE-unix": {
      "fileSize": 4096,
      "isDir": true,
      "modTime": "2018-05-31T10:04:14.894378146-07:00",
      "mode": "dtrwxrwxrwx"
    },
    "/tmp/.ICE-unix/1710": {
      "fileSize": 0,
      "isDir": false,
      "modTime": "2018-05-31T10:04:14.894378146-07:00",
      "mode": "Srwxrwxrwx"
    },
    "/tmp/.ICE-unix/767": {
      "fileSize": 0,
      "isDir": false,
      "modTime": "2018-05-31T09:51:09.455191936-07:00",
      "mode": "Srwxrwxrwx"
    },
    "/tmp/.Test-unix": {
      "fileSize": 4096,
      "isDir": true,
      "modTime": "2018-05-31T09:51:07.611191992-07:00",
      "mode": "dtrwxrwxrwx"
    },
    "/tmp/.X0-lock": {
      "fileSize": 11,
      "isDir": false,
      "modTime": "2018-05-31T10:04:15.326378133-07:00",
      "mode": "-r--r--r--"
    },
    "/tmp/.X1024-lock": {
      "fileSize": 11,
      "isDir": false,
      "modTime": "2018-05-31T09:51:10.571191902-07:00",
      "mode": "-r--r--r--"
    },
    "/tmp/.X11-unix": {
      "fileSize": 4096,
      "isDir": true,
      "modTime": "2018-05-31T10:04:15.326378133-07:00",
      "mode": "dtrwxrwxrwx"
    },
    "/tmp/.X11-unix/X0": {
      "fileSize": 0,
      "isDir": false,
      "modTime": "2018-05-31T10:04:15.326378133-07:00",
      "mode": "Srwxrwxr-x"
    },
    "/tmp/.X11-unix/X1024": {
      "fileSize": 0,
      "isDir": false,
      "modTime": "2018-05-31T09:51:10.575191902-07:00",
      "mode": "Srwxrwxr-x"
    }
  },
  "events": []
}
```
