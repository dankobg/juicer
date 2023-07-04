FROM scratch
COPY juicer-linux-amd64 .
ENTRYPOINT ["/juicer-linux-amd64"]