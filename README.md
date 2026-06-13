# Crescendo
## Descripción general
Crescendo es un servicio de streaming de música que permite a los usuarios reproducir, organizar y gestionar contenido musical de forma sencilla.
## Funcionalidades principales
- Reproducción de canciones
- Creación y administración de playlists
- Búsqueda de canciones, albumes y artistas
## Tecnologías utilizadas
- **Base de datos**: PostgreSQL
- **Cache**: Redis
- **Orquestración**: Docker + Traefik
## Arquitectura
El sistema se basa en una arquitectura simple orientada a servicios:
- PostgreSQL como sistema principal de almacenamiento
- Redis para cacheo de consultas frecuentes
- Docker para contenerización del entorno
- Traefik como proxy reverso y gestor de tráfico

El servicio estará dividido de la siguiente forma:
- Servidor: encargado de proveer la informacion y música al usuario. Internamente, trabajaremos con una arquitectura de microservicios:
    - Backend: trabajará como API entre la base de datos / cache y los usuarios
    - Streaming: manejara la transmisión de la música 
- Cliente: vista y reproduccion de la música / playlists

## Puesta en marcha
Para el uso del servicio se debe seguir los siguientes pasos:
1. Construccion de las imagenes
2. Composición de los contenedores
3. Ejecución del cliente
### Construcción de imágenes
#### API

```bash
cd server/api

docker build -t crescendo-api:1.0.0 .
```
La nomenclatura de la imagen debe seguir el formato:
> ```text
> crescendo-api:<version>
> ```

#### Streaming

```bash
cd server/streaming

docker build -t crescendo-streaming:1.0.0 .
```
La nomenclatura de la imagen debe seguir el formato:
> ```text
> crescendo-streaming:<version>
> ```
### Composición de contenedores

Una vez construidas las imágenes, se deben levantar los servicios definidos en Docker Compose:
```bash
cd Docker

docker compose up -d
```
Para verificar el estado de los contenedores:
```bash
docker compose ps
```
Para detener los servicios:
```bash
docker compose down
```
## Licencia
Este proyecto está licenciado bajo la GNU GPL v3.0.  Para más detalles, ver el archivo `LICENSE` o visitar https://www.gnu.org/licenses/
## Integrantes
- Ignacio Bianchi
- Gianluca Pluchino
- Valentina Huar Lopez 
- Lautaro Rebeco
