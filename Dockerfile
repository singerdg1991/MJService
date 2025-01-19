# Golang backend environment build container
FROM golang:1.23.1 AS build
ENV DEBIAN_FRONTEND noninteractive

WORKDIR /app
ADD . .

ENV GOOS linux
ENV GOARCH amd64

RUN cp -rf ./.gitconfig /root/.gitconfig
RUN go build -o ./build/httpMaja ./cmd/http/*.go
RUN go build -o ./build/majacli ./cmd/cli/*.go
RUN go build -o ./build/openctl ./cmd/openapi/*.go
RUN ./build/openctl --path=public/apidocs --internal-dir-path=internal

# Deploy container
FROM ubuntu AS deploy
ENV DEBIAN_FRONTEND noninteractive

# make sure the path you are calling the httpMaja is the same as .env
WORKDIR /app
COPY --from=build /app/build/httpMaja ./httpMaja
COPY --from=build /app/build/majacli ./majacli
COPY --from=build /app/.env ./.env

COPY --from=build /app/public ./public

EXPOSE 80

# Start The Project
# ENTRYPOINT ["/app/majacli"]

CMD ["/app/httpMaja"]
