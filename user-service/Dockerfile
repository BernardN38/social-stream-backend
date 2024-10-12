
FROM alpine:latest

RUN mkdir /app

COPY userApp /app

RUN apk --no-cache add curl

# RUN apk --no-cache add go
CMD [ "/app/userApp"]