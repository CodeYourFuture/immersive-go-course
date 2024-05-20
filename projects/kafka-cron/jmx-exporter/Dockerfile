FROM openjdk:17
RUN mkdir -p /jmx
COPY ./jmx_prometheus_httpserver-0.17.2.jar /jmx/
COPY ./config.yaml /jmx/
WORKDIR /jmx
ENTRYPOINT ["java","-jar","/jmx/jmx_prometheus_httpserver-0.17.2.jar","8999","/jmx/config.yaml"]