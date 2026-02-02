FROM golang:1.25.0-alpine

WORKDIR /app

# 3. Copiamos los archivos de dependencias primero (optimiza la velocidad)
COPY go.mod go.sum ./
RUN go mod download

# 4. Copiamos el resto del código fuente
COPY . .

# 5. Compilamos la aplicación
RUN go build -o main .

# 6. Comando para iniciar la API
CMD ["go", "run", "main.go"]