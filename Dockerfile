FROM ubuntu
RUN apt update && apt install -y g++ gdb build-essential

ADD main /
ENTRYPOINT [ "/main" ]