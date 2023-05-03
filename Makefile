arm_name = raspifan-arm
remote_path = /home/pi/fan
pi_name = raspi-thai

all: build-arm execute-on-raspi

make-build-folder:
	mkdir -p build

build-amd64: make-build-folder
	go build -o build/raspisensor ./main.go

build-arm: make-build-folder
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(arm_name) ./main.go

execute-on-raspi: create-remote-folder copy-to-raspi copy-config-to-raspi
	# Führt immer alte versionen aus??
	ssh $(pi_name) 'cd $(remote_path) && sudo ./$(arm_name)'

create-remote-folder:
	ssh $(pi_name) mkdir -p $(remote_path)

copy-to-raspi:
	scp build/$(arm_name) $(pi_name):$(remote_path)

copy-config-to-raspi:
	scp config.yml $(pi_name):$(remote_path)