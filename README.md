# PRMV

## Uses

- Golang 1.22

## Requirements

- Docker
- Docker Compose

## Installation

1. Clone this repository:

   ```sh
   git clone https://github.com/rostis232/prmv.git
   cd prmv
   ```

2. Rename the file .env.example to .env

   ```sh
   mv .env.example .env
   ```

3. Fill in the environment variables:
   ```
   PG_PORT=
   PG_PASS=
   PG_USER=
   PG_DB_NAME=
   PORT=
   ```
4. Start Docker Compose:
   ```sh
   docker-compose up -d
   ```
   
   or use Makefile:
   ```sh
   make up
   ```
   
   This will run docker containers with application and PostrgeSQL.

## Usage

The web portal will be available once Docker Compose is up and running.

## Migrations

App uses [golang-migrate](https://github.com/golang-migrate/migrate) for mirgations handling.
Migrations are applied independently when building containers.

## Contact

- Telegram: [rostis232](https://t.me/rostis232)
- Email: [rostislav.pylypiv@gmail.com](mailto:rostislav.pylypiv@gmail.com)
- LinkedIn: [Rostyslav Pylypiv](https://www.linkedin.com/in/rostyslav-pylypiv/)

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
