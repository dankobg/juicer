FROM scratch
WORKDIR /
COPY juicer .
ENTRYPOINT ["/juicer"]
