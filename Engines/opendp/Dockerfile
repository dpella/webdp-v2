FROM python:3.12.2-alpine3.19

WORKDIR /app

RUN apk add --no-cache libgcc --no-cache libc6-compat

RUN pip install \
    dataclasses \
    opendp \
    Flask \
    Flask-RESTful \
    requests \
    pandas


COPY ./src /app/
COPY ./README.md /app/


EXPOSE 8000

CMD ["python", "app.py"]