ARG PORT_STL_BACKEND=8080
ARG MONGODB_URI=mongodb://stluser:stluser@localhost:27017/stl?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+1.3.1

FROM golang:alpine3.15

ENV PORT_STL_BACKEND=${PORT_STL_BACKEND}
ENV MONGODB_URI=${MONGODB_URI}

RUN apk update
RUN apk add bash

RUN mkdir -p /stl/backend
COPY go.mod /stl/
COPY backend/*.go /stl/backend/

RUN cd /stl/backend && \
  go get -u -v && \
  go build -o stl-backend .

WORKDIR /stl/backend
ENTRYPOINT ["./stl-backend"]