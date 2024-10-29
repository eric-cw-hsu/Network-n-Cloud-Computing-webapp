set -e

# environment variable
export DEBIAN_FRONTEND=noninteractive
export CHECKPOINT_DISABLE=1
export TZ=AMERICA/LOS_ANGELES

sudo apt update
sudo apt upgrade -y
sudo apt install nginx -y
sudo apt clean -y