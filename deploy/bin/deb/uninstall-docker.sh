#!/bin/bash

set -e

FORCE=0

# Check if -f (force) is passed
if [[ "$1" == "-f" ]]; then
  FORCE=1
fi

confirm_uninstall() {
  echo "⚠️ This will uninstall Docker Engine and all related components and images, volumes, etc."
  read -p "Are you sure you want to continue? (y/N): " confirm
  if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Uninstallation cancelled."
    exit 0
  fi
}

uninstall_docker() {
  echo "🧹 Removing Docker-related packages..."
  #sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras || true
  sudo dpkg -r docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras || true

  echo "🧹 Deleting Docker data directories (/var/lib/docker and /var/lib/containerd)..."
  sudo rm -rf /var/lib/docker
  sudo rm -rf /var/lib/containerd

  echo "🧹 Removing Docker APT source and keyring..."
  sudo rm -f /etc/apt/sources.list.d/docker.list
  sudo rm -f /etc/apt/keyrings/docker.asc
}

# Main logic
if [[ "$FORCE" -eq 0 ]]; then
  confirm_uninstall
fi

uninstall_docker

echo "✅ Docker has been successfully uninstalled."
