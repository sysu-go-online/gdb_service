FROM ubuntu
RUN apt update && apt install -y g++ gdb build-essentialexit

ADD main /
ENTRYPOINT [ "/main" ]