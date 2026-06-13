## Requisitos

- Docker en ejecución (`docker ps`)
- Contenedores levantados (`docker compose up -d`)
- Archivo `.env` ubicado en `docker/.env`

---

## Ejecución general

Todos los comandos deben ejecutarse desde la carpeta `docker/`:

```bash
cd docker
```
### PostgreSQL
#### Backup
```bash
./scripts/backup-postgres.sh
```
Salida:
```bash
backups/postgres_YYYYMMDD_HHMMSS.sql.gz
```
#### Restore
```bash
./scripts/restore-postgres.sh backups/postgres_xxx.sql.gz
```
### MongoDB

#### Backup
```bash
./scripts/backup-mongo.sh
```
Salida:
```bash
backups/mongo_YYYYMMDD_HHMMSS.archive.gz
```
#### Restore
```bash
./scripts/restore-mongo.sh backups/mongo_xxx.archive.gz
```
## Windows
Se requiere uno de los siguientes programas para ejecutar
### Windows Subsystem for Linux (WSL2)
```bash
wsl
cd /mnt/c/ruta/al/proyecto/docker
./scripts/backup-postgres.sh
```
### Git Bash
```
bash scripts/backup-postgres.sh
```
## Linux
``` bash
cd docker
./scripts/backup-postgres.sh
```
Puede que requiera permisos extra
```bash
chmod +x scripts/*.sh
```
## Notas importantes
- Los backups se almacenan en `docker/backups/`
- Las variables de entorno se cargan desde `docker/.env`
- Los contenedores deben estar en ejecución (`docker ps`)