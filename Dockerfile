FROM golang
WORKDIR /go/src/github.com/vlam321/Inf191BloomFilter
COPY . /go/src/github.com/vlam321/Inf191BloomFilter
ARG service
ARG shard
ARG port
ENV SERVICE=$service
ENV SHARD=$shard
EXPOSE $port
CMD ./run.sh $SERVICE

