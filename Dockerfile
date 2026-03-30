FROM golang:1.23.6 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o /out/app .

FROM oraclelinux:9
#RUN dnf -y install dnf && dnf clean all
COPY --from=build /out/app /usr/local/bin/app
#ENTRYPOINT ["/usr/local/bin/app"]
ENTRYPOINT ["sleep","infinity"]