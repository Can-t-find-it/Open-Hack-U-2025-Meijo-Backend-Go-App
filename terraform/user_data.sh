#!/bin/bash

# ログを /var/log/user_data.log に保存（デバッグ用）
exec > >(tee /var/log/user_data.log|logger -t user-data -s 2>/dev/console) 2>&1

echo "Start Setup..."

#-----------------------------------
# 1. スワップ領域の確保 (t2.micro対策)
#-----------------------------------
# メモリ不足でビルドが落ちないように2GBのスワップファイルを作成
dd if=/dev/zero of=/swapfile bs=128M count=16
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
# 再起動後も有効になるように設定
echo "/swapfile swap swap defaults 0 0" >> /etc/fstab

echo "Swap created."

#-----------------------------------
# 2. OS更新 & 基本ツール
#-----------------------------------
dnf update -y
dnf install -y git docker

#-----------------------------------
# 3. Dockerプラグインのインストール
#-----------------------------------
mkdir -p /usr/local/lib/docker/cli-plugins/

# (A) Docker Compose (V2) のインストール
curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /usr/local/lib/docker/cli-plugins/docker-compose
chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

# (B) Docker Buildx のインストール (←今回追加した部分)
# ※手動で成功したバージョン(v0.17.1)を指定しています
curl -L https://github.com/docker/buildx/releases/download/v0.17.1/buildx-v0.17.1.linux-amd64 -o /usr/local/lib/docker/cli-plugins/docker-buildx
chmod +x /usr/local/lib/docker/cli-plugins/docker-buildx

#-----------------------------------
# 4. Docker起動設定
#-----------------------------------
systemctl start docker
systemctl enable docker
usermod -aG docker ec2-user

#-----------------------------------
# 5. アプリケーションのデプロイ
#-----------------------------------
cd /home/ec2-user

# リポジトリをクローン（URLは適宜ご自身のものか確認してください）
git clone https://github.com/Can-t-find-it/Open-Hack-U-2025-Meijo-Backend-Go-App.git

# 権限修正
chown -R ec2-user:ec2-user /home/ec2-user/Open-Hack-U-2025-Meijo-Backend-Go-App

# ディレクトリ移動
cd /home/ec2-user/Open-Hack-U-2025-Meijo-Backend-Go-App

# .env ファイルの作成 (空ファイルを作成してエラー回避)
touch .env
chown ec2-user:ec2-user .env

#-----------------------------------
# 6. コンテナ起動
#-----------------------------------
echo "Starting Docker Compose..."
# Docker Compose でビルド＆起動
docker compose up -d --build

echo "Setup All Completed!"