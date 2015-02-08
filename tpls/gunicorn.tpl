FROM {{.Registry}}/python:super

WORKDIR {{.AppDir}}

ADD dir/gunicorn/* {{.AppDir}}/
ADD {{.Tarball}} {{.AppDir}}/

RUN [[ -f requirements.txt ]] && pip install -r requirements.txt || echo "no requirements.txt"
RUN [[ -f pip_requirements.txt ]] && pip install -r pip_requirements.txt || echo "no pip_requirements.txt"

EXPOSE 8080

RUN chmod +x control
CMD ["./control", "start", "8080"]

