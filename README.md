# tgbot-ytdl

tgbot-ytdl – телеграм бот для скачивания видео с ютуба, написанный на Go

(над неймингом стоило получше поработать)

---
## Технологии

- Go 1.25.5
- Telegram Bot API (Go lib)
- ytdl в формате гошной библиотеки 
- Какой-то редактор ID3 тегов
---

## Запуск
### Linux/MacOS
1. Клонировать репозиторий:
   ```bash
   https://github.com/diemensa/tgbot-ytdl-go
   cd tgbot-ytdl-go
   
2. Установить переменную окружения BOT_TOKEN
   ```bash
   export BOT_TOKEN="PUT_HERE_BOTFATHER_TOKEN"   

3. Запустить бота:
   ```bash
   go run main.go
