FROM {{.Registry}}/nodejs:base

WORKDIR {{.AppDir}}

ADD dir/nodejs/* {{.AppDir}}/
ADD {{.Tarball}} {{.AppDir}}/

EXPOSE 8080

RUN chmod +x control
CMD ["./control", "start", "8080"]

