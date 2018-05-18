TARGET=31fcn

all: drawin linux win

linux: 
	GOOS=linux GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src

drawin:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src

win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/${TARGET}.exe ./src
	GOOS=windows GOARCH=386 go build -o ./bin/${TARGET}-i386.exe ./src

clean:
	rm -rf ./bin/${TARGET}_*
