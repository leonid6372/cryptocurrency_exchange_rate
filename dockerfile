FROM golang:1.23.0

# Установка wait-for-it для ожидания поднятия сервисов
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

COPY . /crypto_ex_rate

WORKDIR /crypto_ex_rate

RUN go get ./cmd/crypto_ex_rate

CMD /usr/local/bin/wait-for-it crypto_ex_rate_db:5435 --timeout=30 -- \
go run ./cmd/crypto_ex_rate/main.go & \
wait