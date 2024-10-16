# Курсовой проект 14-го потока route256 go middle

- [Домашние задания](./docs)


##  Полезные команды которые запускаются в корне проекта:
- `make run-all` - запуск приложения
- `make .proto-generate` - генерация кода из proto файлов
- `make .serve-swagger` - запуск swagger
- `make .install-goose` - установка goose

## Полезные команды которые запускаются внутри сервисов:
- `make test-cover` - запуск тестов с покрытием
- `make test-e2e` - запуск e2e тестов
- `make lint` - запуск линтера
- `goose create <name_of_migration> sql` - создание миграции
- `minimock -i <InterfaceName> -o <dir/output_mock.go>` - генерация моков
