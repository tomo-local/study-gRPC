#!/bin/bash
set -e

# 引数でユーザー名とメールを受け取る
GIT_USER_NAME=${1:-'github-actions[bot]'}
GIT_USER_EMAIL=${2:-'github-actions[bot]@users.noreply.github.com'}

git config --global user.name "$GIT_USER_NAME"
git config --global user.email "$GIT_USER_EMAIL"
git add .
git diff --cached --quiet || git commit -m "Auto-generate pb.go files from proto"
git push
