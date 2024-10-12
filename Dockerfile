# Etapa 1: Builder - Compilar o binário
FROM golang:latest AS builder

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia apenas os arquivos go.mod e go.sum para aproveitar cache das dependências
COPY go.mod go.sum ./

# Baixa as dependências necessárias
RUN go mod download

# Copia o restante dos arquivos do projeto
COPY . .

# Compila o binário do projeto sem dependências de glibc (usando CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rater_limit ./cmd/rater_limit

# Etapa 2: Final - Criação da imagem mínima para execução
FROM debian:bullseye-slim

# Define o diretório de trabalho na imagem final
WORKDIR /app

# Copia o binário compilado da etapa anterior
COPY --from=builder /app/rater_limit .

# Copia o arquivo init.sql (caso seja necessário para inicializar o banco de dados)
COPY cmd/rater_limit/init.sql .

# Exponha as portas necessárias (ajuste conforme necessidade)
EXPOSE 8080

# Comando de inicialização do contêiner
CMD ["./rater_limit"]
