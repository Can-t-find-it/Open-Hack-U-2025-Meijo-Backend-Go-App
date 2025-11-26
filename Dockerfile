# ---------------------------------
# ステージ1: ビルド環境 (Builder)
# ---------------------------------
FROM golang:1.25-alpine AS builder

# ★★★ 修正点 ★★★
# プロジェクトのファイルに触れる前に、Air をインストールする
# これにより、go.mod の影響を完全に回避します
RUN GOWORK=off go install github.com/air-verse/air@latest

# 作業ディレクトリを作成・設定
WORKDIR /app

# 先に go.mod と go.sum だけをコピーする
COPY go.mod go.sum ./

# Goの依存関係（Ginなど）をダウンロードする
RUN go mod download

# プロジェクトの全ソースコードを作業ディレクトリにコピーする
COPY . .

# (これは本番用のビルドコマンドとして残しておく)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/server

# ---------------------------------
# ステージ: 開発環境 (Development)
# ---------------------------------
# "builder" ステージを引き継いで "development" ステージを定義
FROM builder AS development

# Air は "builder" ステージで既にインストール済みなので、
# ここで再度インストールする必要はありません。

# このコンテナが起動した時のデフォルトコマンドを "air" に設定
CMD ["air"]

# ---------------------------------
# ステージ2: 実行環境 (Final)
# (これは本番イメージを作るためにそのまま残しておきます)
# ---------------------------------
FROM alpine:latest  
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]