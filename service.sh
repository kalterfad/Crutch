#!/bin/bash

user=$(whoami)

echo "[Unit]
Description=Crutch service
After=network.target

[Service]
ExecStart=/usr/local/bin/service-crutch
Restart=always
User=$user

[Install]
WantedBy=multi-user.target" > ./crutch.service
