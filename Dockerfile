FROM ubuntu:14.04
EXPOSE 8000
RUN mkdir /opt/app
COPY ./main /opt/app/main
COPY ./templates /opt/app/templates
WORKDIR /opt/app
CMD ["/opt/app/main"]