#Запуск:

sudo docker build -t db . && sudo docker run -p 5000:5000 -p 5431:5432 -p 5005:5005 db

Порт для прокси: 5005, порт для api: 5000.
