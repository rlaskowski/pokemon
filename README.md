## Pokemon

### Struktura projektu
- error.go - lista błedów do zwrócenia przez serwere
- service.go - startowanie usług systemowych tj. HttpServer 
- router.go - definicja endpointów
- pokemon.go - pobieranie danych z zewnętrznego serwisu oraz logika biznesowa
- http_*.go - definicja serwera http 

### Flagi
- HttpPort - port serwera http
- ServiceAddr - adres usługi zewnętrznej pobierającej dane 