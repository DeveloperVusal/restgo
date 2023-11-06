# Introdution
This project making for the [search engine](https://butago.com/) and also for my portfolio.

For first, execute command:
```bash
npm install
```

Next execute this command for run server:
```bash
npm run serve
```

# For migrate databases
For using migrates of database need install utility the [golang-migrate](https://github.com/golang-migrate/migrate)

## Instruction of installation

### Install on Windows
First need install the scoop
```bash
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```
Next install golang-migrate
```bash
scoop install migrate
```

### Install on Linux
```bash
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
```

### Install on Mac
```bash
brew install golang-migrate
```

## Instruction of migrate

So that will create migrations, run one of the command below:
```bash
make migrate
make migrate dbext=json
```
> The command `make migrate` as default creating sql files