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
## Licencia
Este proyecto está licenciado bajo la GNU GPL v3.0.  Para más detalles, ver el archivo `LICENSE` o visitar https://www.gnu.org/licenses/
## Integrantes
- Ignacio Bianchi
- Gianluca Pluchino
- Valentina Huar Lopez 
- Lautaro Rebeco
