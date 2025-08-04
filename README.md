# Go Webservice mit Logging und Tracking

Dies ist ein Webservice in Go, der einfache mathematische Operationen bereitstellt.

## Funktionen

* **5 Routen:** Addition, Subtraktion, Multiplikation, Division und letzte Ergebnisse.
* **Logging:** Verwendet `zerolog` für schnelles und strukturiertes Logging.
* **Speicherung:** Die Ergebnisse werden standardmäßig im Arbeitsspeicher gespeichert. Über die Umgebungsvariable `USE_DB=true` kann auf eine MySQL-Datenbank umgeschaltet werden.
* **Tracking:** Verwendet `influxdb-client-go` um die Ausführungsdauer der Operationen an eine InfluxDB-Instanz zu senden.
* **Docker Compose:** Eine `docker-compose.yml` Datei zum einfachen Starten aller Dienste (Webservice, MySQL, InfluxDB).
* **Makefile:** Enthält nützliche Befehle zum Builden und Ausführen der Anwendung.

## Konfiguration

Die Anwendung kann über die Datei `config/config.yaml` konfiguriert werden. Alternativ können auch Umgebungsvariablen verwendet werden, die Vorrang vor den Werten in der Konfigurationsdatei haben.

**Konfigurationsoptionen:**

```yaml
database:
  enabled: false  # entspricht USE_DB
  dsn: ""        # entspricht MYSQL_DSN

influxdb:
  url: ""        # entspricht INFLUX_URL
  token: ""      # entspricht INFLUX_TOKEN
  org: ""        # entspricht INFLUX_ORG
  bucket: ""     # entspricht INFLUX_BUCKET

server:
  port: 8080     # Port, auf dem der Webservice läuft
```

## Starten der Anwendung

1.  **Voraussetzungen:** Docker und Docker Compose müssen installiert sein.
2.  **Abhängigkeiten:** Stelle sicher, dass die Go-Abhängigkeiten installiert sind: `go mod tidy`.
3.  **Docker Compose starten:**
    ```sh
    docker-compose up --build
    ```
4.  **Datenbank einrichten (optional, falls `USE_DB=true`):**
    Wenn du die Datenbank verwenden möchtest, musst du die Tabelle erstellen:
    ```sh
    make setup-db
    ```

## Verwendung

Der Webservice läuft standardmäßig auf `http://localhost:8080`.

**Beispielanfragen:**

* **Addition:** 
```bash
curl http://localhost:8080/add?SummandOne=10&SummandTwo=5
```


* **Subtraktion:** 
```bash
curl http://localhost:8080/subtract?Minuend=20&Subtrahend=8
```
* **Multiplikation:** 

```bash
curl http://localhost:8080/multiply?FaktorOne=7&FaktorTwo=6
```

* **Division:** 

```bash
curl http://localhost:8080/divide?Dividend=100&Divisor=10
```

* **Letzte Ergebnisse:** 

```bash
curl http://localhost:8080/results?RecentN=10
```

Der Rückgabewert ist eine Liste der letzten 5 (oder `RecentN`) Operationen.