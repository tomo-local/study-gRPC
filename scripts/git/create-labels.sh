#!/usr/bin/env bash
# GitHub PRに必要なラベルを作成するスクリプト

set -e

services="$1"

if [ -z "$services" ]; then
  echo "No services provided"
  exit 0
fi

# 各サービス用のラベルを作成（存在しない場合のみ）
for service in $services; do
  # ランダムカラー生成
  color=$(printf "%06x" $((RANDOM * RANDOM % 16777216)))

  # サービス別の説明
  description="$service micro service"

  # ラベルが存在するかチェック
  if ! gh label list --search "$service" --limit 1 | grep -q "$service"; then
    echo "Creating label: $service with color: $color"
    gh label create "$service" --description "$description" --color "$color"
  else
    echo "Label $service already exists, skipping creation"
  fi
done
