.PHONY: build run clean setup-db

# Build-Befehl
build:
	go build -o webservice ./cmd/webservice

# Docker Compose starten
run:
	docker-compose up --build

# Docker Compose herunterfahren und Volumes löschen
clean:
	docker-compose down -v

# Datenbank-Initialisierung (für die MySQL-Tabelle)
setup-db:
	docker-compose exec mysql mysql -uuser -ppassword webservice -e "CREATE TABLE IF NOT EXISTS results (id INT AUTO_INCREMENT PRIMARY KEY, operation VARCHAR(50), input_a FLOAT, input_b FLOAT, output FLOAT, timestamp DATETIME);"