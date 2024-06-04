FROM python:3.10.13-slim-bookworm

RUN apt-get update && \
    apt-get install -y default-jre-headless

#WORKDIR /app

# RUN apt-get update && \
#     apt-get install -y wget && \
#     wget -O postgresql.jar https://jdbc.postgresql.org/download/postgresql-42.7.3.jar && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*


RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY requirements.txt /usr/src/app/

RUN pip3 install --no-cache-dir -r requirements.txt

COPY . /usr/src/app

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
