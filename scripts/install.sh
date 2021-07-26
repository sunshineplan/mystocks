#! /bin/bash

OS=$(uname)
if [ $OS = "Linux" ]; then
    ON_LINUX=1
    INSTALL_PATH=/var/www/mystocks
    SOCKET=/run/mystocks.sock
    LOG=/var/log/app/mystocks.log
elif [ $OS = "Darwin" ]; then
    INSTALL_PATH=$HOME/www/mystocks
    SOCKET=none
    LOG=none
else
    abort "Only supported on macOS and Linux."
fi

installSoftware() {
    if [ ${ON_LINUX-} ]; then
        apt -qq -y install nginx mongodb-org-tools
    else
        [ -x $(command -v brew) ] || abort "Require Homebrew."
        brew tap mongodb/brew
        brew install go node mongodb/brew/mongodb-database-tools
    fi
}

installMyStocks() {
    if [ ${ON_LINUX-} ]; then
        mkdir -p $INSTALL_PATH
        curl -Lo- https://github.com/sunshineplan/mystocks/releases/latest/download/release-linux.tar.gz | tar zxC $INSTALL_PATH
        chmod +x $INSTALL_PATH/mystocks
    else
        TMPDIR=$(mktemp -d)
        curl -Lo- https://github.com/sunshineplan/mystocks/archive/main.tar.gz | tar zxC $TMPDIR
        cd $TMPDIR/*
        go build -ldflags "-s -w" && npm i && npm run build || exit 1
        mkdir -p $INSTALL_PATH
        cp -r public mystocks config.ini.default $INSTALL_PATH
        rm -rf $TMPDIR
    fi
}

configMyStocks() {
    read -p 'Please enter metadata server: ' server
    read -p 'Please enter VerifyHeader header: ' header
    read -p 'Please enter VerifyHeader value: ' value
    while true
    do
        read -p 'Use Universal ID(default: false): ' universal
        [ -z $universal ] && universal=false && break
        [ $universal = true -o $universal = false ] && break
        echo Use Universal ID must be true or false!
    done
    read -p "Please enter unix socket(default: $SOCKET): " unix
    [ -z $unix ] && unix=$SOCKET
    [ $unix = none ] && unix=
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter refresh time(second, default: 3): ' refresh
    [ -z $refresh ] && refresh=3
    read -p "Please enter log path(default: $LOG): " log
    [ -z $log ] && log=$LOG
    [ $log = none ] && log=
    read -p 'Please enter update URL: ' update
    read -p 'Please enter exclude files: ' exclude
    [ $log ] && mkdir -p $(dirname $log)
    sed "s,\$server,$server," $INSTALL_PATH/config.ini.default > $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$header/$header/" $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$value/$value/" $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$universal/$universal/" $INSTALL_PATH/config.ini
    sed -i"" -e "s,\$unix,$unix," $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$host/$host/" $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$port/$port/" $INSTALL_PATH/config.ini
    sed -i"" -e "s/\$refresh/$refresh/" $INSTALL_PATH/config.ini
    sed -i"" -e "s,\$log,$log," $INSTALL_PATH/config.ini
    sed -i"" -e "s,\$update,$update," $INSTALL_PATH/config.ini
    sed -i"" -e "s|\$exclude|$exclude|" $INSTALL_PATH/config.ini
    $INSTALL_PATH/mystocks install || exit 1
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

setupNGINX() {
    cp -s $INSTALL_PATH/scripts/mystocks.conf /etc/nginx/conf.d
    sed -i"" -e "s/\$domain/$domain/" $INSTALL_PATH/scripts/mystocks.conf
    sed -i"" -e "s,\$unix,$unix," $INSTALL_PATH/scripts/mystocks.conf
    service nginx reload
}

main() {
    installSoftware
    installMyStocks
    configMyStocks
    if [ ${ON_LINUX-} ]; then
        read -p 'Please enter domain:' domain
        writeLogrotateScrip
        setupNGINX
        service mystocks start
    fi
}

main
