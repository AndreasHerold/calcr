.PHONY: build run clean setup-db

# Build-Befehl
build:
	go build -o webservice ./cmd/webservice

# podman Compose starten
run:
	podman-compose up --build

# podman Compose herunterfahren und Volumes löschen
clean:
	podman-compose down -v

# Datenbank-Initialisierung (für die MySQL-Tabelle)
setup-db:
	podman-compose exec mysql mysql -uuser -ppassword calcr -e ""