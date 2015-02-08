FROM {{.Registry}}/nginxphp:super

RUN mkdir -p /opt/t

ADD {{.Tarball}} /opt/t/

RUN [[ -f /opt/t/config/nginx.conf ]] && mv /opt/t/config/nginx.conf /opt/nginx/conf/nginx.conf || echo "no nginx.conf"
RUN [[ -f /opt/t/config/php.ini ]] && mv /opt/t/config/php.ini /opt/php/etc/php.ini || echo "no php.ini"
RUN [[ -f /opt/t/config/php-fpm.conf ]] && mv /opt/t/config/php-fpm.conf /opt/php/etc/php-fpm.conf || echo "no php-fpm.conf"

RUN rm -rf {{.AppDir}}
RUN mv /opt/t/htdocs {{.AppDir}}

RUN mkdir -p /opt/bin
ADD dir/php/* /opt/bin/

WORKDIR {{.AppDir}}

EXPOSE 8080

RUN chmod +x /opt/bin/control
CMD ["/opt/bin/control", "start", "8080"]

