#Documentación de las Rutas

## 1. Rutas para Administrador:

- **GET /api/admin/books**: Recupera libros filtrados por estado y editorial.
- **GET /api/admin/books/unavailable**: Recupera libros que no están disponibles.
- **GET /api/admin/users**: Recupera los usuarios filtrados por tipo de socio.
- **GET /api/admin/loans**: Obtiene los préstamos activos de los usuarios, filtrados por un rango de fechas.
- **GET /api/admin/fines**: Obtiene las multas pendientes de pago por parte de los usuarios.
- **GET /api/admin/fines**: Obtiene el historial completo de multas de un usuario.
- **GET /api/admin/reservations**: Recupera las reservas activas, filtradas por usuario o libro.
- **GET /api/admin/loans/history**: Obtiene el historial completo de préstamos de un usuario.
- **GET /api/admin/books/available**: Recupera los libros disponibles, filtrados por género y autor.
- **GET /api/admin/books/published**: Recupera los libros ordenados por fecha de publicación.

## 2. Rutas para Usuario:

- **GET /api/books**: Obtiene los libros disponibles filtrados por género y autor (sin acceso administrativo).
- **GET /api/loans**: Obtiene el estado de los préstamos activos del usuario.
- **GET /api/loans/completed**: Obtiene el historial de préstamos completados del usuario.
- **GET /api/fines**: Obtiene las multas pendientes del usuario.
- **GET /api/reservations**: Recupera las reservas activas del usuario.

## 3. Rutas Adicionales:

- **POST /api/login**: Permite a un usuario o administrador iniciar sesión.
- **POST /api/register**: Permite registrar un nuevo usuario.
- **PUT /api/admin/users/{id}**: Permite actualizar los detalles de un usuario específico.
- **DELETE /api/admin/users/{id}**: Elimina un usuario de la base de datos.
- **POST /api/admin/books**: Crea un nuevo libro en la base de datos.
- **PUT /api/admin/books/{id}**: Actualiza la información de un libro.
- **DELETE /api/admin/books/{id}**: Elimina un libro de la base de datos.
- **DELETE /api/reservations/{id}**: Cancela una reserva activa de un usuario.
- **POST /api/loans/extend/{id}**: Extiende un préstamo específico de un usuario.

