FROM golang:alpine

WORKDIR /app

RUN apk add build-base python3-dev && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools wheel

COPY . .
RUN go install
RUN make install

CMD ["discord-tofu"]
