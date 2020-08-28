# golang + opencv + gocv 样例项目

- 基础库制作: `docker build -f Dockerfile-base -t ohko/gocv-base-440 .`

- 项目制作: `docker build -f Dockerfile -t ohko/opencv_test .`

- 项目测试: `docker run --rm -it ohko/opencv_test`