SrcDir=/src
SeetaFaceEngineCommit=0f73c0964cf229d16fe584db14c08c61b1d84105
SeetaFaceEngineDetectionDir=$(SrcDir)/SeetaFaceEngine/FaceDetection
SeetaFaceEngineDetectionBuildDir=$(SeetaFaceEngineDetectionDir)/build
SeetaFaceEngineDetectionSo=libseeta_facedet_lib.so
SeetaFaceEngineDetectionSoPath=$(SeetaFaceEngineDetectionBuildDir)/$(SeetaFaceEngineDetectionSo)
SeetaFaceEngineCMake=cmake3

# clone seeta source
seeta/src/c:
	mkdir -p $(SrcDir)
	cd $(SrcDir) && git clone https://github.com/seetaface/SeetaFaceEngine.git
	cd $(SrcDir)/SeetaFaceEngine && git checkout $(SeetaFaceEngineCommit)
.PHONY: seeta/src/c

# build seeta source
seeta/src/b:
	rm -rf $(SeetaFaceEngineDetectionBuildDir); mkdir $(SeetaFaceEngineDetectionBuildDir)
	cd $(SeetaFaceEngineDetectionBuildDir); $(SeetaFaceEngineCMake) ..; make -j${nproc}
	cp $(SeetaFaceEngineDetectionSoPath) /lib64/$(SeetaFaceEngineDetectionSo)
.PHONY: seeta/src/b

# clone and build seeta source
seeta/src: seeta/src/c seeta/src/b
.PHONY: seeta/src

# build seeta cmd
seeta/b:
	g++ -std=c++11 ./SeetaFaceDetection.cpp \
		-o $(SeetaFaceEngineDetectionBuildDir)/seeta \
		-I $(SeetaFaceEngineDetectionDir)/include/ \
		-L $(SeetaFaceEngineDetectionBuildDir) \
		-lseeta_facedet_lib -ljsoncpp \
		`pkg-config opencv --cflags --libs`
.PHONY: seeta/b

# clone, build seeta source and build seeta cmd
seeta: seeta/src seeta/b
.PHONY: seeta

# clean
clean:
	rm -f /lib64/$(SeetaFaceEngineDetectionSo)
.PHONY: clean
