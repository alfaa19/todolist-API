# Gunakan image Docker resmi untuk Go sebagai dasar
FROM golang:latest

# Buat direktori kerja
WORKDIR /app

ENV MYSQL_USER=${MYSQL_USER}
ENV MYSQL_PASSWORD=${MYSQL_PASSWORD}
ENV MYSQL_NAME=${MYSQL_NAME}
ENV MYSQL_HOST=${MYSQL_HOST}
ENV MYSQL_PORT=${MYSQL_PORT}



# Salin file go.mod dan go.sum ke direktori kerja
COPY go.mod .
COPY go.sum .

# Download dan instal dependensi
RUN go mod download

# Salin seluruh kode sumber aplikasi ke direktori kerja
COPY . .

# Kompilasi aplikasi
RUN go build -o todolist-api

# Jalankan aplikasi ketika container dimulai
CMD ["/app/todolist-api"]