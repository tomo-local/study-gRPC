#!/usr/bin/env bash
# PRの変更ファイルから対象サービスを検出して、カンマ区切りで出力する

set -e

files="$1"
services=(auth greeting memo note)
found=()

for f in $files; do
  for s in "${services[@]}"; do
    if [[ $f == service/$s/* ]]; then
      if [[ ! " ${found[@]} " =~ " $s " ]]; then
        found+=("$s")
      fi
    fi
  done

done

# 改行区切りで出力
for s in "${found[@]}"; do
  # ラベルが存在しなければ自動追加（色はデフォルト: blue）
  if ! gh label list | grep -q "^$s[[:space:]]"; then
    gh label create "$s" --color "1e90ff" --description "$s service"
  fi
  echo "$s"
done
