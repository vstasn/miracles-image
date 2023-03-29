#!/bin/sh
docker build -t miracles-image .
docker tag miracles-image ghcr.io/vstasn/miracles-image:latest
docker push ghcr.io/vstasn/miracles-image:latest
