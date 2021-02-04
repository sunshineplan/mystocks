#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
}

installMyStocks() {
    mkdir -p /var/www/mystocks
    curl -Lo- https://github.com/sunshineplan/mystocks/releases/download/v1.0/release-linux.tar.gz | tar zxC /var/www/mystocks
    cd /var/www/mystocks
    chmod +x mystocks
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
    read -p 'Please enter unix socket(default: /run/mystocks.sock): ' unix
    [ -z $unix ] && unix=/run/mystocks.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/mystocks.log): ' log
    [ -z $log ] && log=/var/log/app/mystocks.log
    read -p 'Please enter update URL: ' update
    read -p 'Please enter exclude files: ' exclude
    mkdir -p $(dirname $log)
    sed "s,\$server,$server," /var/www/mystocks/config.ini.default > /var/www/mystocks/config.ini
    sed -i "s/\$header/$header/" /var/www/mystocks/config.ini
    sed -i "s/\$value/$value/" /var/www/mystocks/config.ini
    sed -i "s/\$universal/$universal/" /var/www/mystocks/config.ini
    sed -i "s,\$unix,$unix," /var/www/mystocks/config.ini
    sed -i "s,\$log,$log," /var/www/mystocks/config.ini
    sed -i "s/\$host/$host/" /var/www/mystocks/config.ini
    sed -i "s/\$port/$port/" /var/www/mystocks/config.ini
    sed -i "s,\$update,$update," /var/www/mystocks/config.ini
    sed -i "s|\$exclude|$exclude|" /var/www/mystocks/config.ini
    ./mystocks install || exit 1
    service mystocks start
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

createCronTask() {
    cp -s /var/www/mystocks/scripts/stock.cron /etc/cron.monthly/mystocks
    chmod +x /var/www/mystocks/scripts/mystocks.cron
}

setupNGINX() {
    cp -s /var/www/mystocks/scripts/mystocks.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/mystocks/scripts/mystocks.conf
    sed -i "s,\$unix,$unix," /var/www/mystocks/scripts/mystocks.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installMyStocks
    configMyStocks
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main
