FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wget apache2 \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir /opt/backend

RUN cd /opt && wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
#     # PATH=$PATH:/usr/local/go/bin

COPY src/app-frontend/index.html /var/www/html
COPY src/app-backend/app.go /opt/backend

EXPOSE 80

CMD ["apachectl", "-D", "FOREGROUND"]