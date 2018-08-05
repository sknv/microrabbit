#!/usr/bin/env bash
# bootstrap vagrant development environment

# update the system
sudo apt update

# install docker
sudo apt install -y --no-install-recommends \
  apt-transport-https \
  ca-certificates \
  curl \
  software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt install -y --no-install-recommends docker-ce

# allow executing docker without sudo
sudo usermod -aG docker ${USER}

# install docker-compose
sudo curl -L https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m) \
  -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

echo 'All set, rock on!'
