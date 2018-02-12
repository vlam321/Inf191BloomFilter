FROM golang
WORKDIR /go/src/github.com/vlam321/Inf191BloomFilter
COPY . /go/src/github.com/vlam321/Inf191BloomFilter
ARG service
ENV SERVICE=$service
CMD ./run.sh $SERVICE

