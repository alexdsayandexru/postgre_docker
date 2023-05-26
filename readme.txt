docker-compose up --build

Если не будет подключаться, скорее всего запущен еще один экземпляр
lsof -n -i:5432 | grep LISTEN

go get -u github.com/lib/pq