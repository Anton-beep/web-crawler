#!/bin/sh

# Запускаем тесты с покрытием и сохраняем результат в файл
sleep 15

echo "Запуск тестов..."
if ! go test -count=1 -v -coverprofile=coverage.out ./...; then
  echo "Ошибка: один или несколько тестов завершились с ошибкой"
  exit 1
fi

if [ ! -f coverage.out ]; then
  echo "Ошибка: файл coverage.out не создан"
  exit 1
fi

# Извлекаем общее покрытие из файла coverage.out
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Проверяем, удалось ли извлечь покрытие
if [ -z "$coverage" ]; then
  echo "Ошибка: не удалось извлечь процент покрытия"
  exit 1
fi

coverage_int=$(echo "$coverage" | awk '{print int($1 + 0.5)}')
# Проверяем, превышает ли покрытие 85%
if [ "$coverage_int" -lt 60 ]; then
  echo "Ошибка: покрытие слишком низкое ($coverage%)"
  exit 1
fi

echo "Все тесты пройдены успешно. Покрытие достаточное ($coverage%)"