SrcDir=/src
SeetaCommitID=0f73c0964cf229d16fe584db14c08c61b1d84105
SeetaFDSrcDir=$(SrcDir)/SeetaFaceEngine/FaceDetection
CMake3=cmake3


seeta/src/c:
	mkdir -p $(SrcDir)
	cd $(SrcDir) && git clone https://github.com/seetaface/SeetaFaceEngine.git
	cd $(SrcDir)/SeetaFaceEngine && git checkout $(SeetaCommitID)
.PHONY: seeta/src/c

seeta/src/b:
	rm -rf $(SeetaFDSrcDir)/build; mkdir $(SeetaFDSrcDir)/build
	cd $(SeetaFDSrcDir)/build; $(CMake3) ..; make -j${nproc}
	cp $(SeetaFDSrcDir)/build/libseeta_facedet_lib.so /lib64/libseeta_facedet_lib.so
.PHONY: seeta/src/b

seeta/src: seeta/src/c seeta/src/b
.PHONY: seeta/src

faced:
	cd libfaced && g++ -std=c++11 faced.cpp -fPIC -shared -o libfaced.so `pkg-config opencv --cflags --libs` \
		-I$(SeetaFDSrcDir)/include/ \
		-L$(SeetaFDSrcDir)/build \
		-lseeta_facedet_lib -ljsoncpp
	cd libfaced && g++ -std=c++11 faced_cmd.cpp -o faced \
		-I$(SeetaFDSrcDir)/include \
		-L$(SeetaFDSrcDir)/build \
		-L. -lfaced
	cd libfaced && LD_LIBRARY_PATH=. ./faced ../face.jpg
.PHONY: seeta

goserver/b:
	go build server.go
	LD_LIBRARY_PATH=./libfaced ./server

goserver: goserver/b
.PHONY: goserver

clean:
	rm -f /lib64/libfaced_facedet_lib.so
	rm -f libfaced/faced
	rm -f libfaced/libfaced.so
.PHONY: clean
