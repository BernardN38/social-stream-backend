FROM alpine:latest

RUN mkdir /app

COPY authApp /app

RUN apk --no-cache add curl

# RUN apk --no-cache add go
CMD [ "/app/authApp"]