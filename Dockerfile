FROM ubuntu:latest
LABEL authors="ayaan"

ENTRYPOINT ["top", "-b"]