# Crutch

This program is designed to scan directories and help you manage files. It allows you to constantly monitor the 
appearance of new files and automatically sort them into folders.


## Getting started

### Prerequisites
* Install **[Go](https://go.dev/)**.

### Create it
Run source build:
```sh
make
```
Create and enable the service:
```sh
make enable
```
Next, start the systemd service:
```sh
make start
```

### Usage
After starting the service, add a directory for scanning:
```sh
crutch add /path/to/dir
```
Crutch will offer a choice: use the base rules or create your own.

When creating your own rules, Crutch will open Nano and prompt you to edit the base rules.

After adding a directory to the scan list, Crutch will scan the selected directory for new files and sort them according to the specified rules. Currently, the rules are created in a json file, and in the future, the process of adding rules will be simplified

To disable directory scanning, enter:
```sh
crutch rm /path/to/dir
```
---
You can also use makefile commands to control the service.

Service stop:
```sh
make stop
```
Service restart:
```sh
make restart
```
Service status:
```shell
make status
```
Removing and disabling the service:
```sh
make disable
```
Or use systemctl commands using the service name "crutch.service"

### Uninstall
To completely uninstall the application, do the following:
```sh
make clean
```