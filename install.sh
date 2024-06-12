#! /usr/bin/env sh

go build .

if test -f /usr/local/bin/synchronex; then
  sudo rm /usr/local/bin/synchronex
fi
sudo cp ./synchronex /usr/local/bin/synchronex

sudo chown "$USER":"$USER" /usr/local/bin/synchronex
sudo chmod +x /usr/local/bin/synchronex