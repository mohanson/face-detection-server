src_dir=$(shell pwd)
seeta_commit_id=0f73c0964cf229d16fe584db14c08c61b1d84105
seeta_fd_src_dir=$(src_dir)/SeetaFaceEngine/FaceDetection
cmake=cmake

seeta/clone:
	cd $(src_dir) && git clone https://github.com/seetaface/SeetaFaceEngine.git
	cd $(src_dir)/SeetaFaceEngine && git checkout $(seeta_commit_id)

seeta/build:
	rm -rf $(seeta_fd_src_dir)/build
	mkdir $(seeta_fd_src_dir)/build
	cd $(seeta_fd_src_dir)/build && $(cmake) .. && make -j${nproc}

seeta: seeta/clone seeta/build

faced:
	cd libfaced && g++ -std=c++11 faced.cpp -fPIC -shared -o libfaced.so `pkg-config opencv4 --cflags --libs` \
		-I$(seeta_fd_src_dir)/include/ \
		-L$(seeta_fd_src_dir)/build \
		-lseeta_facedet_lib -ljsoncpp
	cd libfaced && g++ -std=c++11 faced_cmd.cpp -o faced \
		-I$(seeta_fd_src_dir)/include \
		-L$(seeta_fd_src_dir)/build -L. \
		-lfaced -lseeta_facedet_lib
	LD_LIBRARY_PATH=$(src_dir)/libfaced:$(seeta_fd_src_dir)/build libfaced/faced face.jpg

goserver:
	go build server.go
	LD_LIBRARY_PATH=$(src_dir)/libfaced:$(seeta_fd_src_dir)/build ./server
