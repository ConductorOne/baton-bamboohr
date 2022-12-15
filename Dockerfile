FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-bamboohr"]
COPY baton-bamboohr /