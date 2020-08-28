FROM ohko/gocv-base-440

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR ${GOPATH}/src/opencv

ADD . ${GOPATH}/src/opencv
RUN go build -o opencv_run -i main.go

CMD ["./opencv_run"]
