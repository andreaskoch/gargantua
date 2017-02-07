build:
	go build -o bin/gargantua

install:
	go install

test:
	go test

crosscompile:
	GOOS=linux GOARCH=amd64 go build -o bin/gargantua_linux_amd64
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/gargantua_linux_arm_5
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/gargantua_linux_arm_6
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/gargantua_linux_arm_7
	GOOS=darwin GOARCH=amd64 go build -o bin/gargantua_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o bin/gargantua_windows_amd64
