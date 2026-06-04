# Funcionalidades

Se registran aquí todas las funcionalidades del sistema, indicando las responsabilidades de cada componente.

---

## Funcionalidades de usuario

### USR-001 Registrar usuario

Permite crear una nueva cuenta para acceder a la aplicación.

#### CLI

Muestra la interfaz de registro y realiza validaciones básicas de entrada.

Campos:

- Usuario:
  - Longitud mínima de 6 caracteres.
  - No debe contener caracteres especiales.

- Contraseña:
  - Longitud mínima de 8 caracteres.
  - Debe contener al menos un número.
  - Debe contener al menos un carácter especial.

- Correo electrónico:
  - Debe tener un formato válido.

- Fecha de nacimiento:
  - Debe ser anterior a la fecha actual.

#### API

Responsabilidades:

- Recibir los datos enviados por el CLI.
- Validar nuevamente todos los campos.
- Verificar que el usuario no exista.
- Verificar que el correo no esté registrado.
- Hashear la contraseña.
- Registrar el usuario en la base de datos relacional.

#### Streaming

No participa en esta funcionalidad.

#### Errores

- Nombre de usuario ya existente.
- Correo electrónico ya registrado.
- Datos inválidos.
- Error interno de base de datos.

### USR-002 Iniciar sesión

Ingreso a la aplicación para acceder a la música y playlists.

#### CLI

- Solicita usuario o correo electrónico.
- Solicita contraseña.
- Envía las credenciales a la API.
- Almacena el token o sesión recibida.
- Muestra mensajes de error en caso de credenciales inválidas.

#### API

- Recibe las credenciales del usuario.
- Valida formato de los datos recibidos.
- Busca el usuario en la base de datos.
- Verifica la contraseña mediante comparación con el hash almacenado.
- Genera un token de autenticación o sesión.
- Devuelve los datos necesarios para la autenticación.

#### Streaming

- No participa directamente.
- Valida el token recibido para permitir el acceso a contenido protegido.

---

## Funcionalidades de música

### MUS-001 Reproducir canción

Permite iniciar la reproducción de una canción del catálogo.

#### CLI

- Solicita la reproducción de una canción mediante su identificador.
- Obtiene los metadatos de la canción desde la API.
- Solicita el flujo de audio al servicio de Streaming.
- Reproduce el audio recibido.
- Muestra información de la canción actual.

#### API

- Proporciona información de la canción solicitada.
- Verifica que la canción exista en el catálogo.
- Devuelve los metadatos necesarios para la reproducción.

#### Streaming

- Valida permisos de acceso.
- Obtiene el archivo de audio desde MongoDB.
- Inicia la transmisión del contenido.
- Gestiona la lectura y envío de datos de audio.

---

### MUS-002 Pausar reproducción

Permite detener temporalmente la reproducción actual.

#### CLI

- Detecta la solicitud de pausa.
- Detiene temporalmente la reproducción local.
- Permite posteriormente reanudar la reproducción.

#### API

- No participa directamente.

#### Streaming

- Mantiene disponible el flujo para una posible reanudación.
- Puede registrar el punto actual de reproducción si se implementa persistencia de estado.

---

### MUS-003 Buscar canción

Permite localizar canciones dentro del catálogo.

#### CLI

- Solicita un criterio de búsqueda.
- Muestra los resultados obtenidos.
- Permite seleccionar una canción para reproducirla o agregarla a una playlist.

#### API

- Recibe el término de búsqueda.
- Busca coincidencias por título, artista, álbum o género.
- Devuelve la lista de resultados encontrados.

#### Streaming

- No participa.

---

## Funcionalidades de playlist

### PLY-001 Crear playlist

Permite crear una nueva playlist personal.

#### CLI

- Solicita el nombre de la playlist.
- Envía la solicitud de creación a la API.
- Muestra el resultado de la operación.

#### API

- Valida el nombre de la playlist.
- Asocia la playlist al usuario autenticado.
- Registra la playlist en la base de datos.
- Devuelve la información de la playlist creada.

#### Streaming

- No participa.

---

### PLY-002 Agregar canción a playlist

Permite incorporar una canción a una playlist existente.

#### CLI

- Permite seleccionar una playlist.
- Permite seleccionar una canción.
- Envía la solicitud a la API.
- Muestra confirmación o errores.

#### API

- Verifica que la playlist exista.
- Verifica que el usuario tenga permisos sobre la playlist.
- Verifica que la canción exista.
- Registra la asociación entre playlist y canción.
- Devuelve el resultado de la operación.

#### Streaming

- No participa.