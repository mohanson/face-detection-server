#include <cstdint>
#include <fstream>
#include <iostream>
#include <string>

#include "opencv2/highgui/highgui.hpp"
#include "opencv2/imgproc/imgproc.hpp"

#include "face_detection.h"
#include "jsoncpp/json/json.h"

using namespace std;

int main(int argc, char** argv) {
  if (argc < 2) {
      cout << "Usage: " << argv[0] << " FILE" << endl;
      return -1;
  }

  const char* img_path = argv[1];
  seeta::FaceDetection detector("/src/SeetaFaceEngine/FaceDetection/model/seeta_fd_frontal_v1.0.bin");

  detector.SetMinFaceSize(40);
  detector.SetScoreThresh(2.f);
  detector.SetImagePyramidScaleFactor(0.8f);
  detector.SetWindowStep(4, 4);

  cv::Mat img = cv::imread(img_path, cv::IMREAD_UNCHANGED);
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

  Json::Value resp;
  resp["size"].append(Json::Value(img_data.width));
  resp["size"].append(Json::Value(img_data.height));

  int32_t num_face = static_cast<int32_t>(faces.size());

  for (int32_t i = 0; i < num_face; i++) {
      Json::Value innerResp;
      innerResp["x"] = Json::Value(faces[i].bbox.x);
      innerResp["y"] = Json::Value(faces[i].bbox.y);
      innerResp["width"] = Json::Value(faces[i].bbox.width);
      innerResp["height"] = Json::Value(faces[i].bbox.height);
      resp["faces"].append(Json::Value(innerResp));
  }

  Json::FastWriter fw;
  cout << fw.write(resp) << endl;
}
