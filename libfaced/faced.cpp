#include <cstdint>
#include <fstream>
#include <iostream>
#include <string>

#include "opencv2/highgui/highgui.hpp"
#include "opencv2/imgproc/imgproc.hpp"
#include "jsoncpp/json/json.h"
#include "face_detection.h"

#include "faced.h"

#include <cstring>
#include <string>

using namespace std;

extern seeta::FaceDetection detector("/src/SeetaFaceEngine/FaceDetection/model/seeta_fd_frontal_v1.0.bin");

const char* FaceDetect(char* path) {
    detector.SetMinFaceSize(40);
    detector.SetScoreThresh(2.f);
    detector.SetImagePyramidScaleFactor(0.8f);
    detector.SetWindowStep(4, 4);

    cv::Mat img = cv::imread(path, cv::IMREAD_UNCHANGED);
    cv::Mat img_gray;

    if (img.channels() != 1)
      cv::cvtColor(img, img_gray, cv::COLOR_BGR2GRAY);
    else
      img_gray = img;

    seeta::ImageData img_data;
    img_data.data = img_gray.data;
    img_data.width = img_gray.cols;
    img_data.height = img_gray.rows;
    img_data.num_channels = 1;

    std::vector<seeta::FaceInfo> faces = detector.Detect(img_data);

    Json::Value root;
    root["size"].append(Json::Value(img_data.width));
    root["size"].append(Json::Value(img_data.height));

    int32_t num_face = static_cast<int32_t>(faces.size());

    for (int32_t i = 0; i < num_face; i++) {
        Json::Value innerResp;
        innerResp["x"] = Json::Value(faces[i].bbox.x);
        innerResp["y"] = Json::Value(faces[i].bbox.y);
        innerResp["width"] = Json::Value(faces[i].bbox.width);
        innerResp["height"] = Json::Value(faces[i].bbox.height);
        root["face"].append(Json::Value(innerResp));
    }

    char *out = new char[root.toStyledString().length() + 1];
    std::strcpy(out, root.toStyledString().c_str());
    return (const char*)out;
}
