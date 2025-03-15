BUILD_DIR=build

.PHONY:cmd clean

all:cmd

cmd:
	CXX=${CXX} CC=${CC} go build -o ${BUILD_DIR}/assisant app/assisant.go

clean:
	rm -rf ${BUILD_DIR}/*




