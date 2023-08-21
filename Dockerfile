FROM golang:1.20.3-alpine

#set working directory
WORKDIR /app

# copy over mod and sum
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy the rest of the files
COPY . . 

# run the program
RUN go build -o main .

# expose port 
EXPOSE 3000

# set the entry point of the container to the executable
CMD ["./main"]