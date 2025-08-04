# Go Webservice mit Logging und Tracking

Dies ist ein Webservice in Go, der einfache mathematische Operationen bereitstellt.

## Funktionen

* **5 Routen:** Addition, Subtraktion, Multiplikation, Division und letzte Ergebnisse.
* **Logging:** Verwendet `zerolog` für schnelles und strukturiertes Logging.
* **Speicherung:** Die Ergebnisse werden standardmäßig im Arbeitsspeicher gespeichert. Über die Umgebungsvariable `USE_DB=true` kann auf eine MySQL-Datenbank umgeschaltet werden.
* **Tracking:** Verwendet `influxdb-client-go` um die Ausführungsdauer der Operationen an eine InfluxDB-Instanz zu senden.
* **Docker Compose:** Eine `docker-compose.yml` Datei zum einfachen Starten aller Dienste (Webservice, MySQL, InfluxDB).
* **Makefile:** Enthält nützliche Befehle zum Builden und Ausführen der Anwendung.

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

* **Addition:** `http://localhost:8080/add?SummandOne=10&SummandTwo=5`
* **Subtraktion:** `http://localhost:8080/subtract?SummandOne=20&SummandTwo=8`
* **Multiplikation:** `http://localhost:8080/multiply?SummandOne=7&SummandTwo=6`
* **Division:** `http://localhost:8080/divide?SummandOne=100&SummandTwo=10`
* **Letzte Ergebnisse:** `http://localhost:8080/results?RecentN=10`

Der Rückgabewert ist eine Liste der letzten 5 (oder `RecentN`) Operationen.