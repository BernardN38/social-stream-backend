FROM alpine:latest

RUN mkdir /app

COPY postApp /app

RUN apk --no-cache add curl

# RUN apk --no-cache add go
CMD [ "/app/postApp"]