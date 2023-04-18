.SILENT:

all: program

program: dependencies client service install

dependencies: go.mod
	go mod download

client: cmd/client/crutch.go
	go build -o crutch ./cmd/client/crutch.go
	echo "build client"

service: cmd/service/serviceCrutch.go
	go build -o service-crutch ./cmd/service/serviceCrutch.go
	echo "build service"

install: client service
	if [ ! -d "$(HOME)/.crutch" ]; then \
    	mkdir $(HOME)/.crutch && cp rules.json $(HOME)/.crutch/; \
    fi
	sudo mv crutch /usr/local/bin/
	sudo mv service-crutch /usr/local/bin/
	echo "move client, service, create .crutch"

.PHONY: systemd
systemd: service
	sudo cp crutch.service /etc/systemd/system/
	sudo systemctl enable crutch.service

.PHONY: start
start:
	sudo systemctl start crutch.service

.PHONY: stop
stop:
	sudo systemctl stop crutch.service

.PHONY: restart
restart:
	sudo systemctl restart crutch.service

.PHONY: clean
clean:
	sudo rm -rf /usr/local/bin/crutch
	sudo rm -rf /usr/local/bin/service-crutch
	sudo rm -rf /etc/systemd/system/crutch.service
	sudo rm -rf $(HOME)/.crutch
