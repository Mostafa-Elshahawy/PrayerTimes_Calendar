FROM golang:1.21.6-alpine AS build-env

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM public.ecr.aws/lambda/provided:al2023

WORKDIR /app

COPY --from=build-env /app/main /app/.env ./

CMD ["./main"]
