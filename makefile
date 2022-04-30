DST=.
BIN=genWaymark.exe

all:
	go build -o $(DST)/$(BIN) -trimpath

release:
	go build -x -o $(DST)/$(BIN) -trimpath -a -tags netgo -installsuffix netgo -ldflags="-s -w -extldflags=\"static\" -buildid="
	./upx $(DST)/$(BIN) --lzma -o $(DST)/$(BIN).upx
	rm $(DST)/$(BIN)
	mv $(DST)/$(BIN).upx $(DST)/$(BIN)
