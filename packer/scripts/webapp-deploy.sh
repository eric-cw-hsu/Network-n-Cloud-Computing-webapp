# setup webapp (go)

set -e

sudo mkdir -p /opt/webapp
sudo mkdir -p /opt/bak-webapp

sudo groupadd csye6225
sudo useradd --system -g csye6225 -s /usr/sbin/nologin csye6225

# setup nginx-default
sudo rm /etc/nginx/sites-enabled/default
sudo cp /tmp/nginx.conf /etc/nginx/sites-available/webapp.conf
sudo ln -s /etc/nginx/sites-available/webapp.conf /etc/nginx/sites-enabled/
sed -i 's/worker_connections 768;/worker_connections 1024;/g' /etc/nginx/nginx.conf
sed -i 's/#\s*multi_accept on;/multi_accept on;/g' /etc/nginx/nginx.conf

# main server
sudo cp /tmp/app /opt/webapp/app
sudo cp -r /tmp/migrations /opt/webapp/migrations
sudo cp /tmp/app.service /etc/systemd/system/app.service

# bak server
sudo cp /tmp/app /opt/bak-webapp/app
sudo cp -r /tmp/migrations /opt/bak-webapp/migrations
sudo cp /tmp/app-bak.service /etc/systemd/system/app-bak.service

sudo systemctl daemon-reload

sudo systemctl enable app
sudo systemctl enable app-bak
sudo systemctl enable nginx

sudo chown -R csye6225:csye6225 /opt/webapp
