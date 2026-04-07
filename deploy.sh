#!/bin/bash
set -e

# ===========================================
# Deploy script для origo_api
# ===========================================
# Первая установка:  bash deploy.sh init
# Обновление:        bash deploy.sh update
# Остановка:         bash deploy.sh down
# Логи:              bash deploy.sh logs
# ===========================================

COMPOSE_FILE="docker-compose.prod.yml"

case "$1" in
  init)
    echo "🚀 Первый запуск origo_api..."

    if [ ! -f .env ]; then
      echo "❌ Файл .env не найден!"
      echo "   cp .env.prod.example .env && nano .env"
      exit 1
    fi

    # Проверяем что origo-admin сеть существует
    if ! docker network ls | grep -q "origo-admin_origo-network"; then
      echo "❌ Сеть origo-admin_origo-network не найдена!"
      echo "   Сначала запустите origo-admin: cd /opt/origo-admin && bash deploy.sh init"
      exit 1
    fi

    docker compose -f $COMPOSE_FILE up -d --build
    echo ""
    echo "✅ origo_api запущен на порту 8081!"
    ;;

  update)
    echo "🔄 Обновление origo_api..."
    git pull origin main
    docker compose -f $COMPOSE_FILE build api
    docker compose -f $COMPOSE_FILE up -d
    echo "✅ Обновление завершено!"
    ;;

  update-no-cache)
    echo "🔄 Обновление origo_api..."
    git pull origin main
    docker compose -f $COMPOSE_FILE build --no-cache api
    docker compose -f $COMPOSE_FILE up -d
    echo "✅ Обновление завершено!"
    ;;

  down)
    echo "🛑 Остановка origo_api..."
    docker compose -f $COMPOSE_FILE down
    echo "✅ Остановлено."
    ;;

  logs)
    docker compose -f $COMPOSE_FILE logs -f --tail=100
    ;;

  status)
    docker compose -f $COMPOSE_FILE ps
    ;;

  *)
    echo "Использование: bash deploy.sh {init|update|down|logs|status}"
    exit 1
    ;;
esac
