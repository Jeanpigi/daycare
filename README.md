# Daycare API 🧸🏫

![Go Lint & Build](https://github.com/Jeanpigi/daycare/actions/workflows/lint.yml/badge.svg)
![Security Scan](https://github.com/Jeanpigi/daycare/actions/workflows/trivy.yml/badge.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)


Backend REST API para la gestión de una guardería, desarrollado en Go (Golang) con MySQL, orientado a registrar niños, controlar asistencias (ingreso y salida), calcular cobros y administrar precios y promociones.

## 📌 Descripción general

### Este proyecto implementa el backend de un sistema de guardería donde:

1) Un administrador configura precios, promociones y usuarios.
2) El personal (staff) registra niños y controla su asistencia diaria.
3) El sistema calcula automáticamente el valor a cobrar según:
4) tiempo de permanencia
5) precio base activo
6) promociones vigentes

Todas las acciones administrativas quedan registradas para auditoría. El backend está diseñado para ser consumido por un frontend (por ejemplo Vue.js) u otros clientes (Postman, curl, apps móviles, etc.).

## 📌 Características principales

- Autenticación con **JWT**
- Roles de usuario (`ADMIN`, `STAFF`)
- Registro de niños por documento
- Control de **check-in / check-out**
- Cálculo automático de cobros
- Sistema de **precios configurables**
- **Promociones** por tiempo o días acumulados
- Auditoría de acciones administrativas
- Arquitectura limpia y escalable

## 🧱 Arquitectura del proyecto

El proyecto sigue una arquitectura por capas clara y mantenible:

- cmd/api
Punto de entrada de la aplicación. Aquí se inicializa todo: configuración, base de datos, repositorios, servicios, handlers, middlewares y servidor HTTP.

- internal/config
Carga la configuración desde variables de entorno y construye el DSN de la base de datos.

- internal/db
Maneja la conexión a MySQL y el registro del driver.

- internal/domain
Define los modelos del negocio (User, Child, Attendance, Pricing, Promotion, etc.).
No contiene lógica ni dependencias externas.

- internal/repository/mysql
Acceso a datos. Cada repositorio encapsula las consultas SQL de una tabla o conjunto de tablas.

- internal/service

## Lógica de negocio:

- autenticación
- cálculo de precios
- check-in / check-out
- administración de precios, promociones y usuarios
- bootstrap del primer administrador
- internal/httpapi

## Capa HTTP:

- handlers: endpoints
- middleware: autenticación y roles
- router: definición de rutas
- migrations:Scripts SQL para crear las tablas de la base de datos.


## 🗄️ Base de datos

### Tablas principales

| Tabla | Descripción |
|-----|-----------|
| `users` | Usuarios del sistema (ADMIN / STAFF) |
| `children` | Niños registrados |
| `attendances` | Asistencias (ingreso / salida) |
| `settings_pricing` | Precio base activo |
| `promotions` | Promociones configurables |
| `audit_log` | Auditoría de acciones administrativas |


## ⚙️ Requisitos

- Go **1.21+**
- MySQL / MariaDB **8+**
- Podman o Docker (opcional)
- curl (para pruebas)

## 🚀 Instalación y ejecución

Clonar el repositorio y entrar al directorio del proyecto.

```bash
git clone https://github.com/Jeanpigi/daycare.git
cd daycare

En otra terminal se debe correr

podman-compose up -d build

## Verificar la API
curl http://localhost:8080/health

## Respuesta esperada 
{"ok":true}


## Detener el servicio o contenedor
podman-compose down
podman volume rm daycare_mysql_data



## Configuración

### Definir las variables de entorno necesarias:

- APP_ENV: Entorno de ejecución (dev, prod, etc.)
- HTTP_ADDR: Dirección y puerto donde escucha la API (por ejemplo :8080)
- DB_HOST: Host de la base de datos
- DB_PORT: Puerto de la base de datos
- DB_NAME: Nombre de la base de datos
- DB_USER: Usuario de la base de datos
- DB_PASS: Contraseña del usuario
- JWT_SECRET: Clave secreta para firmar los tokens JWT
- JWT_TTL_MINUTES: Tiempo de vida del token en minutos

### Base de datos (opcional con Podman)

Puedes levantar MySQL usando Podman o Docker. Una vez levantada la base de datos, ejecuta los scripts SQL del directorio migrations en orden para crear las tablas.


## Ejecutar la aplicación

### Instalar dependencias y ejecutar el servidor:

```bash
  go mod tidy
  go run ./cmd/api
```

Verificar que la API esté funcionando accediendo al endpoint de salud:

```bash
curl http://localhost:8080/health
```

Debe responder con un JSON indicando que el servicio está activo.

{"ok":true}

## 🔐 Autenticación y roles

La API utiliza JWT para autenticación.

El token debe enviarse en el header:

Authorization: Bearer <token>

## Roles

- ADMIN

Tiene acceso a:

1) creación de usuarios
2) configuración de precios
3) creación y activación de promociones

- STAFF

Tiene acceso a:

1) registro de niños
2) check-in
3) check-out

## 🧭 Cómo usar la API

### Flujo normal de uso

- Crear el primer usuario administrador (bootstrap).
- Iniciar sesión como administrador.
- Crear usuarios STAFF.
- Registrar niños.
- Registrar ingreso (check-in).
- Registrar salida (check-out).
- El sistema calcula automáticamente el cobro.

## ENDPOINTS:

### NOTA: Este endpoint solo funciona si no existe ningún administrador.

### Crear Usuario administrador:

- Ruta: /admin/bootstrap

Método: POST
Se envía el nombre, correo y contraseña del administrador inicial.

### Iniciar sesión

- Ruta: /auth/login


Método: POST
Se envía el correo y la contraseña. La respuesta incluye un token JWT.

### Crear usuarios STAFF

- Ruta: /admin/users

  
Método: POST
Requiere rol ADMIN.
Permite crear usuarios que registran asistencias.

### Registrar un niño

- Ruta: /children

Método: POST
Se registra el niño con su documento, nombre y datos del acudiente.

### Check-in

- Ruta: /attendances/check-in

Método: POST
Se registra la hora de ingreso del niño usando su número de documento.

### Check-out

- Ruta: /attendances/check-out
  
Método: POST
Se registra la hora de salida.


## 🧾 Auditoría (audit_log)

### La tabla audit_log registra acciones administrativas importantes, como:

- cambios de precios
- creación o activación de promociones
- acciones realizadas por administradores

 ### Esto permite:

- trazabilidad
- control interno
- respaldo ante reclamos

No se usa para operaciones diarias como check-in o check-out.

## 🧠 Decisiones de diseño

1) Separación clara de responsabilidades por capas
2) Lógica de negocio aislada de HTTP y SQL
3) Repositorios enfocados solo en datos
4) Seguridad basada en JWT
5) Pensado para integrarse fácilmente con un frontend (Vue.js u otro)

## 🛠️ Posibles mejoras futuras

- Documentación OpenAPI / Swagger
- Tests unitarios y de integración
- Panel administrativo
- Reportes mensuales
- Despliegue en VPS o cloud


## 👤 Autor

Jean Pierre Giovanni Arenas Ortiz


