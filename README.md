# Starbucks Directory Watcher --> NRI Integration
* Reports contents of a specified directory to NRI Inventory

### [Download The Latest Release HERE](https://source.datanerd.us/xxxxxxxx)

### Dir Watcher NRI Integration Screenshots:
![alt text](https://source.datanerd.us/FIT/DirWatcher/blob/master/images/DirWatcher.jpg "Logo Title Text 1")


### Requirements
* NRI Agent


### NRI Configuration


* Install 3 files noted below in the 3 seperate folders
* `dirwatcher.exe` Executable located in Releases
* Edit `newrelic-infra-exam-config.yml` to add `Dpass`
* Resart NRI Agent
* Verify data in NewRelic under NRI Inventory

### File Structure _(Note 3 different folders & Executable may have a different file extension based on Operating System)_
* `C:\Program Files\New Relic\newrelic-infra\integrations.d\newrelic-infra-dirwatcher-config.yml`
* `C:\Program Files\New Relic\newrelic-infra\custom-integrations\newrelic-infra-dirwatcher-definition.yml`
* `C:\Program Files\New Relic\newrelic-infra\custom-integrations\newrelic-infra-dirwatcher\dirwatcher.exe`


#### Example newrelic-infra-dirwatcher-definition.yml
* Interval currently set to run at 1 minute.
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
#### Example newrelic-infra-dirwatcher-config.yml
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
``` 

### Testing Locally

* You can run the NRI Integration Exectuable locally by manually passing in the required arguments.
* This can help when testing a new feature, or prior to running as an NRI Integration on a new server.

#### Example Command:
```
PS C:\Program Files\New Relic\newrelic-infra\custom-integrations\newrelic-infra-dirwatcher> .\dirwatcher.exe -dir_name C:\crdata\mutex -pretty -inventory
```
### OUTPUT
* You can ignore this WARN message when testing locally as it is expected
```
WARN[0000] Environment variable NRIA_CACHE_PATH is not set, using default C:\Users\ayork\AppData\Local\Temp\exam.exe.json
```
#### Expected JSON Blob
* Use this to confirm you're receiving the expected data
```
{"name":"com.newrelic.dirwatcher","protocol_version":"1","integration_version":"0.1.0","metrics":[{"ActivationHotButtonsDlg.flg":"ENABLED","AllowCreditTransactionVoid.flg":"ENABLED","CupFundRoundUp.flg":"ENABLED","DigitalTips.flg":"ENABLED","FilterAskMe.flg":"ENABLED","LabelItemsCountOn.flg":"ENABLED","LibraOn.xml":"ENABLED","McmUseTls.flg":"ENABLED","NewRefundPolicyOn.flg":"ENABLED","NoRewardsAuthenticate.flg":"ENABLED","event_type":"DirWatcher"}],"inventory":{"ActivationHotButtonsDlg.flg":{"value":"enabled"},"AllowCreditTransactionVoid.flg":{"value":"enabled"},"CupFundRoundUp.flg":{"value":"enabled"},"DigitalTips.flg":{"value":"enabled"},"FilterAskMe.flg":{"value":"enabled"},"LabelItemsCountOn.flg":{"value":"enabled"},"LibraOn.xml":{"value":"enabled"},"McmUseTls.flg":{"value":"enabled"},"NewRefundPolicyOn.flg":{"value":"enabled"},"NoRewardsAuthenticate.flg":{"value":"enabled"}},"events":[]}

```
