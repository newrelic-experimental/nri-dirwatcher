#!/bin/sh

set -e

if [ -z "$1" ]; then
  printf "\nIntegration Name Required:\n"
  printf "\tUsage: \t./build.sh <integration-name>\n\n"
  exit 1
fi

NRI_NAME="nri-$1"
HOME=$(pwd)
PATH_WINDOWS="build/windows/$NRI_NAME"
PATH_LINUX="build/linux/$NRI_NAME"

rm -rf build/ &>/dev/null

if ls *windows.yml &>/dev/null ; then

  echo "Making windows 386 build"
  GOOS=windows GOARCH=386 go build -o "$PATH_WINDOWS/${NRI_NAME}_windows_386.exe" src/*.go

  echo "Making windows amd64 build"
  GOOS=windows GOARCH=amd64 go build -o "$PATH_WINDOWS/${NRI_NAME}_windows_amd64.exe" src/*.go

  echo "Making windows .zip package"za
  #cp etc/install.ps1 $PATH_WINDOWS
  cp *windows.yml $PATH_WINDOWS
  cp README.md $PATH_WINDOWS
  cp LICENSE $PATH_WINDOWS

  echo "\$IntegrationName = \"$NRI_NAME\"" > $PATH_WINDOWS/install.ps1
  cat <<'EOF' >> $PATH_WINDOWS/install.ps1

  $ARCH = $ENV:PROCESSOR_ARCHITECTURE.toLower()

  if (-NOT ($ARCH -eq "amd64")) {
   $ARCH = "x86"
  }

  $ExecutableName = $IntegrationName + "_windows_" + $ARCH + ".exe"

  write-host "`n ## New Relic $IntegrationName Installer ## `n"

  $definition_dir = "C:\Program Files\New Relic\newrelic-infra\custom-integrations"
  $config_dir = "C:\Program Files\New Relic\newrelic-infra\integrations.d"
  $script_dir = Split-Path $script:MyInvocation.MyCommand.Path

  write-host "`n----------------------------"
  write-host " Admin requirement check...  "
  write-host "----------------------------`n"

  ### require admin rights
  if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
     write-Warning "This setup needs admin permissions. Please run this file as admin."
     break
  }

  write-host " ...passed!  "

  write-host "`n----------------------------"
  write-host " Copying files...  "
  write-host "----------------------------`n"

  if (-not (test-path "$definition_dir\$IntegrationName")) {
    New-item $definition_dir\$IntegrationName -itemtype directory
  }

  Copy-item -Force -Recurse $script_dir\$ExecutableName -Destination $definition_dir\$IntegrationName\"$IntegrationName.exe"

  Copy-item -Force $script_dir\$IntegrationName-definition-windows.yml -Destination $definition_dir

  Copy-item -Force $script_dir\$IntegrationName-config-windows.yml -Destination $config_dir


  write-host " ...finished.  "

  write-host "`n---------------------------------------"
  write-host " Restarting New Relic Infrastructure agent... "
  write-host "-----------------------------------------`n"
  $serviceName = 'newrelic-infra'
  $nrServiceInfo = Get-Service -name $serviceName

  if ($nrServiceInfo.Status -ne 'Running')
  {
    write-host "New Relic Infrastructure not running currently, starting..."
    Start-Service -name $serviceName
  } Else {
    Stop-Service -name $serviceName
    Start-Service -name $serviceName
    write-host " Restart complete! "
  }

  Read-Host "`nPress Enter to exit..." | Out-Null

EOF

  cd build/windows
  zip $NRI_NAME.zip "$NRI_NAME"/*
  mv $NRI_NAME.zip ../
  cd $HOME

fi

if ls *linux.yml &>/dev/null; then

  echo "Making linux 386 build"
  GOOS=linux GOARCH=386 go build -o "$PATH_LINUX/${NRI_NAME}_linux_386" src/*.go

  echo "Making linux amd64 build"
  GOOS=linux GOARCH=amd64 go build -o "$PATH_LINUX/${NRI_NAME}_linux_amd64" src/*.go

  echo "Making linux package"
  cp *linux.yml $PATH_LINUX
  cp README.md $PATH_LINUX
  cp LICENSE $PATH_LINUX

  echo "#!/bin/bash" > $PATH_LINUX/install.sh
  echo "INTEGRATION_NAME=\"$NRI_NAME\"" >> $PATH_LINUX/install.sh
  echo "$TEMPLATE_LINUX" >> $PATH_LINUX/install.sh


  cat <<'EOF' >> $PATH_LINUX/install.sh

  printf "Installing $INTEGRATION_NAME Extension...\n"

  DEFINITION_PATH="/var/db/newrelic-infra/custom-integrations"
  CONFIG_PATH="/etc/newrelic-infra/integrations.d"
  SERVICE='newrelic-infra'
  ARCH=`uname -m`

  if [ "$ARCH" = "x86_64" ]; then
    BINARY="${INTEGRATION_NAME}_linux_amd64"
  elif [ "$ARCH" = "i386" ]; then
    BINARY="${INTEGRATION_NAME}_linux_386"
  fi

  if [ -z "$BINARY" ]; then
    printf "\nUnable to determine architecture (x86_64 | i386) via uname -m\n\n"
    exit 1
  fi

  #check os release
  if [ -f /etc/os-release ]; then #amazon/redhat/fedora check
    . /etc/os-release
    OS=$NAME
    VERSION=$VERSION_ID
    printf "OS Name: $OS\n"
    printf "Version: $VERSION\n"
  elif [ -f /etc/lsb-release ]; then #ubuntu/debian check
    . /etc/lsb-release
    OS=$DISTRIB_ID
    VERSION=$DISTRIB_RELEASE
    printf "OS Name: $OS\n"
    printf "Version: $VERSION\n"
  fi

  #check init system
  initCmd=`ps -p 1 | grep init | awk '{print $4}'`
  if [ "$initCmd" == "init" ]; then
    SYSCMD='upstart'
  fi
  sysdCmd=`ps -p 1 | grep systemd | awk '{print $4}'`
  if [ "$sysdCmd" == "systemd" ]; then
    SYSCMD='systemd'
  fi

  printf "Copying files...\n"

  [ -d $DEFINITION_PATH/$INTEGRATION_NAME ] || mkdir -p $DEFINITION_PATH/$INTEGRATION_NAME

  cp *config-linux.yml $CONFIG_PATH/
  cp *definition-linux.yml $DEFINITION_PATH/
  cp $BINARY $DEFINITION_PATH/$INTEGRATION_NAME/$INTEGRATION_NAME


  printf "Script complete. Restarting Infrastructure Agent...\n"
  if [ $SYSCMD == 'systemd' ]; then
    systemctl restart $SERVICE
  elif [ $SYSCMD == 'upstart' ]; then
    initctl restart $SERVICE
  fi
EOF

  chmod +x $PATH_LINUX/install.sh

  cd build/linux
  tar -czvf $NRI_NAME.tar.gz "$NRI_NAME"/*
  mv $NRI_NAME.tar.gz ../
  cd $HOME

fi
