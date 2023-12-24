FROM debian
ADD ./ /app 
WORKDIR /app
EXPOSE 9090
ENTRYPOINT ["/app/main"]