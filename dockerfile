FROM golang:1.23.1 as build

# diretory of work in container
WORKDIR /app

# copy archives go.mod and go.sum to the container
COPY go.mod go.sum ./

# dowload dependencys
RUN go mod dowload

# copy to font-code to the container
COPY . .

# Construct the exec
RUN go build -o main .

# Exec
FROM golang:1.23.1

# define directory of work to execution
WORKDIR /app

# copy the exec of build 
COPY --from=build /app/main .

# the port exposed 
EXPOSE 8080

# execute the app
CMD ["./main"]

