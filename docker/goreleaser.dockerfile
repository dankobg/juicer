FROM scratch
WORKDIR /
COPY juicer-linux-amd64 .
ENTRYPOINT ["/juicer-linux-amd64"]
