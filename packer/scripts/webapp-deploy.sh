# setup webapp (go)

set -e

sudo mkdir -p /opt/webapp

sudo groupadd csye6225
sudo useradd --system -g csye6225 -s /usr/sbin/nologin csye6225

# setup nginx-default
sudo rm /etc/nginx/sites-enabled/default
sudo cp /tmp/nginx.conf /etc/nginx/sites-available/webapp.conf
sudo ln -s /etc/nginx/sites-available/webapp.conf /etc/nginx/sites-enabled/

sudo cp /tmp/app /opt/webapp/app

sudo cp /tmp/app.service /etc/systemd/system/app.service

sudo cp -r /tmp/migrations /opt/webapp/migrations

sudo systemctl daemon-reload

sudo systemctl enable app
sudo systemctl enable nginx

sudo chown -R csye6225:csye6225 /opt/webapp
