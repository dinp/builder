FROM {{.Registry}}/tomcat:7.0.56

RUN mkdir -p /opt/bin
ADD dir/tomcat7/* /opt/bin/

ADD {{.Tarball}} {{.AppDir}}/

WORKDIR {{.AppDir}}
RUN unzip -q {{.Tarball}} && rm -rf {{.Tarball}}

EXPOSE 8080

RUN chmod +x /opt/bin/control
CMD ["/opt/bin/control", "start", "8080"]

