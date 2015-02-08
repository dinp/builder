FROM {{.Registry}}/javaweb:super

RUN mkdir -p /opt/bin
ADD dir/tomcat6/* /opt/bin/

ADD {{.Tarball}} {{.AppDir}}/

WORKDIR {{.AppDir}}
RUN unzip -q {{.Tarball}} && rm -rf {{.Tarball}}

EXPOSE 8080

RUN chmod +x /opt/bin/control
CMD ["/opt/bin/control", "start", "8080"]

