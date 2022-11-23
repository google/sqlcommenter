# First pull Golang image
FROM golang:1.19-alpine

ENV APP_NAME sqlcommenter-http

# Set workdir
WORKDIR $GOPATH/src/$APP_NAME

# Copy application data into image
COPY . .
 
# Budild application
RUN CGO_ENABLED=0 go build -v
 
# Expose application port
EXPOSE 8081
 
# Start app
CMD ["./sqlcommenter-http", "--db_engine=pg"]
