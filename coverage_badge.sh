#!/bin/bash

# Verifica se coverage.out existe
if [ ! -f coverage.out ]; then
    echo "Arquivo coverage.out não encontrado!"
    exit 1
fi

# Extrai a cobertura total do arquivo coverage.out
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
COVERAGE_VALUE=${COVERAGE%\%}  # Remove o símbolo de porcentagem

# Define a cor do badge baseado na porcentagem
COLOR="red"
if [ $COVERAGE_VALUE -ge 80 ]; then
    COLOR="green"
elif [ $COVERAGE_VALUE -ge 50 ]; then
    COLOR="yellow"
fi

# Gera o badge usando shields.io e salva como coverage.svg
curl -o coverage.svg "https://img.shields.io/badge/coverage-$COVERAGE_VALUE%25-$COLOR"

echo "Coverage Badge Atualizado!"