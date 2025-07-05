# 📚 Sistema de Gestión de Biblioteca

Sistema completo de gestión de biblioteca con frontend en React y backend en Go, incluyendo gestión de libros, usuarios, préstamos, multas y reservas.

## 🚀 Tecnologías Utilizadas

### Backend
- **Go** - Lenguaje principal
- **MySQL** - Base de datos
- **Gin** - Framework web
- **JWT** - Autenticación
- **Golang-migrate** - Migraciones de base de datos

### Frontend
- **React** - Framework de UI
- **Vite** - Build tool
- **React Router** - Navegación
- **CSS3** - Estilos

## 📋 Prerrequisitos

- **Go** 1.21 o superior
- **Node.js** 18 o superior
- **MySQL** 8.0 o superior
- **Git**

## 🛠️ Instalación y Configuración

### 1. Clonar el Repositorio
```bash
git clone <url-del-repositorio>
cd biblioteca_db
```

### 2. Configurar Base de Datos

#### Crear Base de Datos MySQL
```sql
CREATE DATABASE biblioteca CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### Configuración del Backend
El backend usa parámetros de línea de comandos. Los valores por defecto son:
- **Base de datos**: `root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true`
- **JWT Secret**: `secreto123`
- **Puerto**: `4000`

### 3. Configurar Backend

#### Instalar Dependencias
```bash
cd server
go mod download
```

#### Ejecutar Migraciones
```bash
# Ejecutar todas las migraciones
make migrate-up

# Revertir migraciones (¡CUIDADO!)
make migrate-down
```

#### Ejecutar Backend
```bash
# Desarrollo (con valores por defecto)
go run cmd/api/main.go

# O con parámetros personalizados
go run cmd/api/main.go -dns="root:tu_password@(127.0.0.1:3306)/biblioteca?parseTime=true" -jwt="tu_secret" -port=4000

# O usando Makefile
make run
```

El backend estará disponible en: `http://localhost:4000`

### 4. Configurar Frontend

#### Instalar Dependencias
```bash
cd client
npm install
```

#### Configurar URL del Backend
El frontend busca la URL del backend en la variable de entorno `VITE_URL_BACKEND`. Por defecto usa `http://localhost:4000/v1/api`.

Para cambiar la URL, crear archivo `.env` en `client/`:
```env
VITE_URL_BACKEND=http://localhost:4000/v1/api
```

#### Ejecutar Frontend
```bash
npm run dev
```

El frontend estará disponible en: `http://localhost:5173`

## 📁 Estructura del Proyecto

```
biblioteca_db/
├── server/                 # Backend en Go
│   ├── cmd/
│   │   ├── api/           # Servidor principal
│   │   └── migrate/       # Migraciones
│   ├── internal/          # Lógica interna
│   │   ├── auth/          # Autenticación
│   │   ├── data/          # Modelos y acceso a datos
│   │   ├── mailer/        # Envío de emails
│   │   └── validator/     # Validaciones
│   └── docs/              # Documentación API
├── client/                # Frontend en React
│   ├── src/
│   │   ├── components/    # Componentes React
│   │   ├── services/      # Servicios API
│   │   ├── hooks/         # Hooks personalizados
│   │   └── styles/        # Estilos CSS
│   └── public/            # Archivos públicos
└── README.md
```

## 🔐 Roles y Permisos

### Usuario Normal
- Ver libros disponibles
- Realizar préstamos
- Ver historial de préstamos
- Pagar multas
- Crear reservas

### Administrador
- Gestión completa de libros
- Gestión de usuarios
- Ver todos los préstamos
- Gestión de multas
- Gestión de reservas
- Crear editoriales y autores

## 🚀 Comandos Útiles

### Backend
```bash
cd server

# Ejecutar en desarrollo
go run cmd/api/main.go

# Con parámetros personalizados
go run cmd/api/main.go -dns="root:password@(127.0.0.1:3306)/biblioteca?parseTime=true" -jwt="secret" -port=4000

# Compilar
go build -o app ./cmd/api

# Ejecutar migraciones
make migrate-up
make migrate-down

# Generar documentación API
make docs
```

### Frontend
```bash
cd client

# Desarrollo
npm run dev

# Build para producción
npm run build

# Preview build
npm run preview
```

## 📚 API Endpoints

### Autenticación
- `POST /v1/api/login` - Iniciar sesión
- `POST /v1/api/register` - Registrar usuario

### Libros
- `GET /v1/api/books` - Obtener libros disponibles
- `GET /v1/api/admin/books` - Gestión de libros (admin)

### Préstamos
- `POST /v1/api/loans` - Crear préstamo
- `GET /v1/api/loans` - Ver préstamos del usuario
- `GET /v1/api/admin/loans` - Ver todos los préstamos (admin)

### Multas
- `GET /v1/api/fines` - Ver multas del usuario
- `PUT /v1/api/fines/{id}/pay` - Pagar multa
- `GET /v1/api/admin/fines` - Gestión de multas (admin)

### Reservas
- `POST /v1/api/reservation` - Crear reserva
- `GET /v1/api/reservations` - Ver reservas del usuario
- `GET /v1/api/admin/reservations` - Gestión de reservas (admin)

## 🔧 Configuración de Desarrollo

### Parámetros del Backend
```bash
# Parámetros disponibles
-port=4000                    # Puerto del servidor
-env=development              # Ambiente (development|staging|production)
-dns="root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true"  # Conexión a BD
-jwt=secreto123               # Secret para JWT
-exp=1h                       # Duración del token JWT
-iss=miApp                    # Emisor del token
-limiter-rps=2                # Requests por segundo
-limiter-burst=4              # Burst del rate limiter
-limiter-enabled=true         # Habilitar rate limiter
-cors-trusted-origins="http://localhost:3000"  # Orígenes CORS
```

### Variables de Entorno Frontend
```env
VITE_URL_BACKEND=http://localhost:4000/v1/api
```

## 🐛 Solución de Problemas

### Error de Conexión a Base de Datos
- Verificar que MySQL esté ejecutándose
- Confirmar credenciales en el parámetro `-dns`
- Verificar que la base de datos `biblioteca` exista

### Error de CORS
- Verificar que `-cors-trusted-origins` incluya `http://localhost:5173`
- Reiniciar el backend después de cambios

### Error de Permisos
- Ejecutar migraciones: `make migrate-up`
- Verificar que los permisos estén correctamente asignados

## 📝 Notas de Desarrollo

- El backend usa parámetros de línea de comandos para configuración
- Las migraciones son incrementales y versionadas
- El frontend usa React Router para navegación
- Los estilos están organizados por componentes
- La API sigue convenciones RESTful

## 🤝 Contribución

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.