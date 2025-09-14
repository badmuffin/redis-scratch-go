## Redis from scratch using go

A simplified in-memory database written in GO that mimics some core features of Redis.

### Features

Currently supported commands

- `PING` - responds with `PONG` or the provided argument
- `SET` – store a key-value pair in memory
- `GET` – retrieve the value of a key
- `HSET` – store a field-value pair inside a hash
- `HGET` – retrieve a value from a hash field
- `HGETALL` – retrieve all fields and values from a hash
