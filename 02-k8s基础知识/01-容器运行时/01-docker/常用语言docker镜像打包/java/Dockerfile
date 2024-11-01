FROM maven:3.6.3-adoptopenjdk-8 as target

ENV MAVEN_HOME /usr/share/maven
ENV PATH $MAVEN_HOME/bin:$PATH
COPY settings.xml /usr/share/maven/conf/
WORKDIR /build
COPY pom.xml .
RUN mvn dependency:go-offline  # use docker cache
COPY src/ /build/src/
RUN mvn clean package -Dmaven.test.skip=true


FROM java:8
WORKDIR /app
RUN  rm /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=target /build/target/*.jar  /app/app.jar
EXPOSE 8080
ENTRYPOINT ["java","-Xmx768m","-Xms256m","-Djava.security.egd=file:/dev/./urandom","-jar","/app/app.jar"]

