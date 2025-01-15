# Микросервис отслеживания и хранения истории курсов криптовалют

После запуска сервис позволяет добавть криптовалюту для хранения истории её цены, запроса цены на указанный момент времени и удаления криптовалюты из списка для отслеживания и хранения цены. Сохранённые данные не удаляются из БД и доступны к просмотру даже после прекращения отслеживания криптовалюты и перезапуска сервиса.

Для запуска микросервиса необходимо сделать следующее:

1) Скачать проект. Настроить частоту обновления цены изменив значения параметра secToUpdate внтури файла cmd/crypto_ex_rate/main.go
2) Запустить командную строку и перейти в каталог проекта
3) При первом запуске прописать "docker-compose up --build crypto_ex_rate" (необходима установка Docker) и дождаться запуска контейнеров Docker. При повторных запусках можно пропускать флаг --build
4) Микросервис находится по адресу 0.0.0.0:8085 (используя API описанное в api_swagger.yaml файле можно добавлять, удалять для отслеживания и просматривать цены криптовалют)
5) Для корректного завершения работы микросервиса необходимо прописать "docker-compose down" находясь в окне командной строки
