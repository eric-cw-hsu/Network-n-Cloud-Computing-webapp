# init database
set -e

echo "Setting up database"

sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql

sudo -u postgres psql -c "CREATE USER ubuntu WITH PASSWORD '1qaz@WSX3edc';"
sudo -u postgres psql -c "CREATE DATABASE webapp WITH OWNER ubuntu;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE webapp TO ubuntu;"
