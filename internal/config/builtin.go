package config

const envPrefix = "gift_card"

var builtinConfig = []byte(`http_server:
  address: 0.0.0.0:8080
database:
  driver: mysql
  host: localhost
  port: 3306
  db: gift-card
  user: gift-card-app
  password: password
user:
  secret: example-secret`)
