# 🚀 Backend - Sistema de Biblioteca

Backend desarrollado en Go para el sistema de gestión de biblioteca, con arquitectura modular y API RESTful.

## 🏗️ Arquitectura

### Estructura del Proyecto
```
server/
├── cmd/
│   ├── api/                    # Punto de entrada principal
│   │   ├── main.go            # Servidor HTTP
│   │   ├── routes.go          # Definición de rutas
│   │   ├── middleware.go      # Middlewares personalizados
│   │   ├── authhandler.go     # Autenticación y registro
│   │   ├── bookhandler.go     # Gestión de libros
│   │   ├── loanhandler.go     # Gestión de préstamos
│   │   ├── finehandler.go     # Gestión de multas
│   │   ├── reservationhandler.go # Gestión de reservas
│   │   ├── userhandler.go     # Gestión de usuarios
│   │   ├── editorialhandler.go # Gestión de editoriales
│   │   ├── autorhandler.go    # Gestión de autores
│   │   └── healcheck.go       # Health check
│   └── migrate/               # Herramientas de migración
│       ├── migration.go       # Ejecutor de migraciones
│       └── migrations/        # Archivos SQL de migración
├── internal/
│   ├── auth/                  # Lógica de autenticación
│   │   ├── auth.go           # Autenticación principal
│   │   └── jwt.go            # Manejo de JWT
│   ├── data/                  # Capa de acceso a datos
│   │   ├── models.go         # Estructuras de datos
│   │   ├── user.go           # Operaciones de usuario
│   │   ├── book.go           # Operaciones de libros
│   │   ├── loan.go           # Operaciones de préstamos
│   │   ├── fine.go           # Operaciones de multas
│   │   ├── reservation.go    # Operaciones de reservas
│   │   ├── editorial.go      # Operaciones de editoriales
│   │   ├── author.go         # Operaciones de autores
│   │   └── permissions.go    # Sistema de permisos
│   ├── jsonlog/              # Logging estructurado
│   ├── mailer/               # Sistema de emails
│   └── validator/            # Validaciones
├── docs/                     # Documentación Swagger
├── go.mod                    # Dependencias Go
├── go.sum                    # Checksums de dependencias
└── Makefile                  # Comandos de automatización
```

## �� Configuración

### Parámetros de Línea de Comandos
El backend usa parámetros de línea de comandos para la configuración:

```bash
# Parámetros principales
-port=4000                    # Puerto del servidor HTTP
-env=development              # Ambiente (development|staging|production)
-dns="root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true"  # Conexión a MySQL
-jwt=secreto123               # Secret para JWT
-exp=1h                       # Duración del token JWT
-iss=miApp                    # Emisor del token

# Rate Limiting
-limiter-rps=2                # Requests por segundo
-limiter-burst=4              # Burst del rate limiter
-limiter-enabled=true         # Habilitar rate limiter

# CORS
-cors-trusted-origins="http://localhost:3000"  # Orígenes permitidos

# Base de datos
-db-max-open-conns=25         # Conexiones máximas abiertas
-db-max-idle-conns=25         # Conexiones máximas inactivas
-db-max-idle-time=15m         # Tiempo máximo de inactividad
```

### Valores por Defecto
```bash
# Configuración por defecto
Puerto: 4000
Ambiente: development
Base de datos: root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true
JWT Secret: secreto123
JWT Expiración: 1h
JWT Emisor: miApp
Rate Limiter: 2 RPS, burst 4
CORS: http://localhost:3000
```

### Dependencias Principales
```go
require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/golang-migrate/migrate/v4 v4.16.2
    golang.org/x/time v0.3.0
)
```

## 🚀 Ejecución

### Desarrollo
```bash
# Instalar dependencias
go mod download

# Ejecutar migraciones
make migrate-up

# Ejecutar servidor con valores por defecto
go run cmd/api/main.go

# O con parámetros personalizados
go run cmd/api/main.go -dns="root:password@(127.0.0.1:3306)/biblioteca?parseTime=true" -jwt="mi_secret" -port=4000
```

### Producción
```bash
# Compilar
go build -o app ./cmd/api

# Ejecutar con parámetros de producción
./app -env=production -dns="user:pass@(host:port)/db?parseTime=true" -jwt="production_secret"
```

### Comandos Makefile
```bash
make build      # Compilar aplicación
make run        # Ejecutar aplicación
make clean      # Limpiar archivos compilados
make migrate-up # Ejecutar migraciones
make migrate-down # Revertir migraciones
make docs       # Generar documentación Swagger
```

## 🔐 Sistema de Autenticación

### JWT (JSON Web Tokens)
- **Algoritmo**: HS256
- **Expiración**: Configurable (default: 1h)
- **Claims**: userId, exp, iat, nbf, iss, aud

### Roles y Permisos
```go
// Permisos disponibles
PermissionBooksRead      = "books:read"
PermissionBooksWrite     = "books:write"
PermissionBooksDelete    = "books:delete"
PermissionLoansCreate    = "loans:create"
PermissionLoansView      = "loans:view"
PermissionLoansManage    = "loans:manage"
PermissionFinesRead      = "fines:read"
PermissionFinesCreate    = "fines:create"
PermissionReservationsCreate = "reservations:create"
PermissionReservationsView   = "reservations:view"
```

### Middleware de Autenticación
```go
// Requerir autenticación
app.requireAuthenticatedUser(handler)

// Requerir permiso específico
app.requirePermission(PermissionBooksRead, handler)
```

## 📊 Base de Datos

### Esquema Principal
- **socio**: Usuarios del sistema
- **libro**: Catálogo de libros
- **prestamo**: Registro de préstamos
- **multa**: Sistema de multas
- **reserva**: Sistema de reservas
- **editorial**: Editoriales
- **autor**: Autores
- **permisos**: Sistema de permisos
- **usuario_permisos**: Relación usuario-permiso

### Migraciones
Las migraciones están versionadas y son incrementales:
- `000001_init_database.up.sql` - Esquema inicial
- `000002_init_procedure.up.sql` - Procedimientos almacenados
- `000003_add_permissions.up.sql` - Sistema de permisos
- `000005_add_reservation_to_loan_procedure.up.sql` - Integración reservas-préstamos

## 🔄 API Endpoints

### Autenticación
```
POST   /v1/api/login          # Iniciar sesión
POST   /v1/api/register       # Registrar usuario
GET    /v1/healthcheck        # Health check
```

### Libros (Usuario)
```
GET    /v1/api/books                    # Libros disponibles
GET    /v1/api/books/reservation        # Libros para reserva
```

### Libros (Admin)
```
GET    /v1/api/admin/books              # Todos los libros
GET    /v1/api/admin/books/{id}/edit    # Libro para editar
POST   /v1/api/admin/books              # Crear libro
POST   /v1/api/admin/books/{id}         # Actualizar libro
DELETE /v1/api/admin/books/{id}         # Eliminar libro
```

### Préstamos
```
POST   /v1/api/loans                    # Crear préstamo
GET    /v1/api/loans                    # Préstamos del usuario
GET    /v1/api/loans/completed          # Historial de préstamos
POST   /v1/api/loans/return/{id}        # Devolver préstamo
POST   /v1/api/loans/extend/{id}        # Extender préstamo
POST   /v1/api/admin/loans              # Préstamos activos (admin)
GET    /v1/api/admin/loans/history      # Historial completo (admin)
```

### Multas
```
GET    /v1/api/fines                    # Multas del usuario
PUT    /v1/api/fines/{id}/pay           # Pagar multa
GET    /v1/api/admin/fines/to           # Multas pendientes (admin)
GET    /v1/api/admin/fines              # Multas por usuario (admin)
GET    /v1/api/admin/fines/search       # Buscar multas (admin)
```

### Reservas
```
POST   /v1/api/reservation              # Crear reserva
GET    /v1/api/reservations             # Reservas del usuario
DELETE /v1/api/reservations/{id}        # Cancelar reserva
GET    /v1/api/admin/reservations       # Todas las reservas (admin)
```

### Usuarios (Admin)
```
POST   /v1/api/admin/users              # Usuarios por tipo
GET    /v1/api/admin/users/all          # Todos los usuarios
```

### Editoriales y Autores
```
POST   /v1/api/editoriales              # Crear editorial
GET    /v1/api/editoriales              # Listar editoriales
POST   /v1/api/admin/autores            # Crear autor
GET    /v1/api/autores                  # Listar autores
```

## 🛡️ Seguridad

### Middlewares Implementados
1. **CORS** - Control de acceso entre dominios
2. **Rate Limiting** - Limitación de requests
3. **Authentication** - Verificación de JWT
4. **Authorization** - Verificación de permisos
5. **Panic Recovery** - Manejo de pánicos
6. **Request Logging** - Logging de requests

### Validaciones
- Validación de email
- Validación de contraseñas
- Validación de datos de entrada
- Sanitización de parámetros

## 📝 Logging

### Estructura de Logs
```json
{
  "level": "info",
  "time": "2024-01-15T10:30:00Z",
  "message": "Request processed",
  "method": "GET",
  "url": "/v1/api/books",
  "status": 200,
  "duration": "15ms"
}
```

## 🧪 Testing

### Ejecutar Tests
```bash
# Tests unitarios
go test ./...

# Tests con coverage
go test -cover ./...

# Tests específicos
go test ./internal/data
```

## 📚 Documentación API

### Swagger
La documentación Swagger está disponible en:
```
http://localhost:4000/docs/
```

### Generar Documentación
```bash
make docs
```

## 🔍 Monitoreo

### Health Check
```
GET /v1/healthcheck
```

### Métricas
- Requests por segundo
- Tiempo de respuesta
- Errores por endpoint
- Uso de memoria

## 🚨 Manejo de Errores

### Códigos de Error
- `400` - Bad Request (datos inválidos)
- `401` - Unauthorized (no autenticado)
- `403` - Forbidden (sin permisos)
- `404` - Not Found (recurso no existe)
- `500` - Internal Server Error (error del servidor)

### Respuestas de Error
```json
{
  "error": "Error message",
  "details": {
    "field": "validation error"
  }
}
```

## 🔧 Desarrollo

### Estructura de Handlers
```go
func (app *application) handlerName(w http.ResponseWriter, r *http.Request) {
    // 1. Validar entrada
    // 2. Procesar lógica de negocio
    // 3. Acceder a datos
    // 4. Responder
}
```

### Convenciones
- Handlers en `cmd/api/`
- Modelos en `internal/data/`
- Validaciones en `internal/validator/`
- Autenticación en `internal/auth/`

## 🚀 Despliegue

### Docker
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
EXPOSE 4000
CMD ["./main", "-env=production", "-dns=user:pass@(host:port)/db?parseTime=true", "-jwt=production_secret"]
```

### Parámetros de Producción
```bash
./app -env=production -dns="user:pass@(host:port)/db?parseTime=true" -jwt="production_secret" -cors-trusted-origins="https://yourdomain.com"
```

