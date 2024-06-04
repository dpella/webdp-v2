FROM python:3.10.13-slim-bookworm

RUN apt-get update && \
    apt-get install -y default-jre-headless

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY ./test/test-requirements.txt /usr/src/app/

RUN pip3 install --no-cache-dir -r test-requirements.txt

COPY ../ /usr/src/app

CMD ["python", "test.py"]
