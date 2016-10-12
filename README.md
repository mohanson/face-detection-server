# FaceDetectionServer

基于 [SeetaFace](https://github.com/seetaface/SeetaFaceEngine) 的高性能人脸识别服务, 使用 Golang 与 CPP 混合开发.

![Face](./face_detection.jpg)

# Requirements

please install them _yourself_.

- cmake3
- opencv
- jsoncpp

# Build and Usage

```
make seeta
make faced
make goserver
```

```
$ http POST :8090/detection/upload < ./face.jpg
```

```json
HTTP/1.1 200 OK
Content-Length: 70
Content-Type: application/json
Date: Wed, 12 Oct 2016 02:47:09 GMT

{
    "face": [
        {
            "Y": 167,
            "height": 287,
            "width": 287,
            "x": 103
        }
    ],
    "size": [
        500,
        650
    ]
}
```
