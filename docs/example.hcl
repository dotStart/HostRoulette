bind-address = "0.0.0.0:8090"

search {
  addresses = [
    "http://localhost:9200"
  ]
}

cache {
  server "vagrant" {
    address = "localhost:6379"
  }

  database = 0
}

auth {
  client-id = "abcdef1234"
  client-secret = "abcdef1234"
}
