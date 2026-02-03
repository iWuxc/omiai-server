#!/bin/bash
set -e

echo "Starting setup on OpenCloudOS 8..."

# 1. Install Docker & Docker Compose
echo "Installing Docker..."
yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

systemctl start docker
systemctl enable docker

# 2. Install Nginx
echo "Installing Nginx..."
yum install -y nginx
systemctl start nginx
systemctl enable nginx

# 3. Setup Directories
mkdir -p /data/omiai-server/deploy
mkdir -p /data/omiai-server/runtime/logs

# 4. Setup Log Rotation for Docker containers
echo '{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}' > /etc/docker/daemon.json
systemctl restart docker

echo "Setup completed successfully!"
echo "Please configure Nginx at /etc/nginx/conf.d/ and add GitHub Secrets."
