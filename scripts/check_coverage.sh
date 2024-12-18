#!/bin/sh

rm -rf cmd

echo "Запуск тестов..."
if ! go test -count=1 -v -coverprofile=coverage.profile ./...; then
  echo "Ошибка: один или несколько тестов завершились с ошибкой"
  exit 1
fi

if [ ! -f coverage.profile ]; then
  echo "Ошибка: файл coverage.profile не создан"
  exit 1
fi

# Extracting the total coverage from the file coverage.profile
coverage=$(go tool cover -func=coverage.profile | grep total | awk '{print $3}' | sed 's/%//')

# Checking whether the coating was removed
if [ -z "$coverage" ]; then
  echo "Ошибка: не удалось извлечь процент покрытия"
  exit 1
fi

coverage_int=$(echo "$coverage" | awk '{print int($1 + 0.5)}')
# Checking if coverage exceeds 60%
if [ "$coverage_int" -lt 60 ]; then
  echo "Ошибка: покрытие слишком низкое ($coverage%)"
  exit 1
fi

echo "Все тесты пройдены успешно. Покрытие достаточное ($coverage%)"

go tool cover -func coverage.profile