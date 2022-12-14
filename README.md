# Запуск тестов:  
- make test  

# Запуск приложения (makefile):  
- сборка и запуск: make init,  
- остановка: make stop,  
- перезапуск: make run.  

# Запуск приложения (docker-compose):  
- сборка и запуск:  
    + docker-compose -f docker-compose.yml down -v --remove-orphans  
    + docker-compose -f docker-compose.yml pull  
    + docker-compose -f docker-compose.yml build --pull  
    + docker-compose -f docker-compose.yml up -d  
- остановка:  
    + docker-compose -f docker-compose.yml down  
- перезапуск: 
    + docker-compose -f docker-compose.yml up -d  

# WEB интерфейс: 
- http://localhost или http://{удаленный ip}    

# Требования:  
[x] 1	Приложение должно иметь веб-интерфейс с отображением десяти последних по времени публикаций.  
[x] 2	Приложение должно принимать на вход конфигурационный файл в формате JSON с массивом ссылок на RSS-ленты информационных сайтов и периодом опроса в минутах.  
[x] 3	Приложение должно регулярно выполнять обход всех переданных в конфигурации RSS-лент.  
[x] 4	Приложение должно выполнять чтение каждой RSS-ленты в отдельном потоке выполнения (горутине).  
[x] 5	Приложение должно сохранять публикации в БД.  
[x] 6	Приложение должно состоять из сервера приложений, базы данных и веб-интерфейса пользователя.  
[x] 7	Веб-интерфейс должен получать от сервера приложений данные в формате JSON.  
[x] 8	Сервер приложения должен предоставлять API, посредством которого осуществляется взаимодействие сервера и веб-интерфейса.  
[x] 9	API должен предоставлять метод для получения заданного количества новостей. Требуемое количество публикаций указывается в пути запроса метода API.  
[x] 10	Агрегатор должен хранить как минимум следующий набор данных для каждой публикации:  
- заголовок,  
- описание,  
- дата публикации,  
- ссылка на источник.  

# Критерии оценки:  
[] 1. Веб-приложение пользователя работоспособно. Отображаются последние новости из источников, указанных в конфигурации.	10  
[] 2. Структура пакетов логична и отражает структуру приложения.	5  
[] 3. Для всех пакетов (кроме, возможно, исполняемого) написаны тесты с достаточным покрытием.	10  
[] 4. Модель данных соответствует требованиям и условиям задачи:  
- структура БД логична,  
- XML с потоком RSS декодируется верно. 10  

[] 5. В каталоге сервера существует корректный файл конфигурации config.json.	2  
[] 6. Присутствует файл со схемой БД (в случае использования реляционной СУБД).	3  
[] 7. Экспортируемые методы снабжены комментариями.	2  
[] 8. Для обхода лент RSS используются отдельные горутины.	3  
[] 9. Для обработки результатов обхода RSS и ошибок используются каналы.	3              