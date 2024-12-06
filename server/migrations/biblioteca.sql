CREATE DATABASE biblioteca;

USE biblioteca;

CREATE TABLE editorial (  
  ideditorial INT PRIMARY KEY,  
  nombre VARCHAR(255),  
  direccion VARCHAR(255),  
  paginaweb VARCHAR(100)  
);  
  
CREATE TABLE libro (  
  idlibro INT PRIMARY KEY,  
  ideditorial INT,  
  fechapublicacion DATE,  
  titulo VARCHAR(255),  
  genero VARCHAR(100),  
  estado ENUM('disponible', 'prestado', 'reservado'),  
  FOREIGN KEY (ideditorial) REFERENCES editorial(ideditorial) 
);  
  
CREATE TABLE autor (  
  idautor INT PRIMARY KEY,  
  nombre VARCHAR(255),  
  nacionalidad VARCHAR(100)  
);  
  
CREATE TABLE socio (  
  idsocio INT PRIMARY KEY,  
  nombre VARCHAR(255),  
  direccion VARCHAR(255),  
  telefono VARCHAR(50),  
  correo VARCHAR(255),  
  fechanacimiento DATE,  
  tiposocio ENUM('normal', 'estudiante', 'profesor'),  
  fecharegistro DATE,  
  imagenperfil VARCHAR(255),  
  rol ENUM('usuario', 'administrador') DEFAULT 'usuario'  
);  
  
CREATE TABLE usuariopassword (   
  idusuario INT PRIMARY KEY,   
  hash_contrasena VARCHAR(255),   
  FOREIGN KEY (idusuario) REFERENCES socio(idsocio)  
);  
  
CREATE TABLE prestamo (  
  idprestamo INT PRIMARY KEY,  
  idsocio INT,  
  idlibro INT,  
  fechaprestamo DATE,  
  fechadevolucion DATE,  
  estado ENUM('activo', 'completado'),  
  FOREIGN KEY (idsocio) REFERENCES socio(idsocio),  
  FOREIGN KEY (idlibro) REFERENCES libro(idlibro)  
);  
  
CREATE TABLE multa (  
  idmulta INT PRIMARY KEY,  
  idprestamo INT,  
  saldopagar DECIMAL(10, 2),  
  fechamulta DATE,  
  estado ENUM('pagada', 'pendiente'),  
  FOREIGN KEY (idprestamo) REFERENCES prestamo(idprestamo)  
);  
  
CREATE TABLE libro_autor (  
  idlibro INT,  
  idautor INT,  
  PRIMARY KEY (idlibro, idautor),  
  FOREIGN KEY (idlibro) REFERENCES libro(idlibro),  
  FOREIGN KEY (idautor) REFERENCES autor(idautor)  
);  
  
CREATE TABLE reserva (  
  idreserva INT PRIMARY KEY,  
  idsocio INT,  
  idlibro INT,  
  fechareserva DATE,  
  estado ENUM('activa', 'cancelada', 'completada'),  
  FOREIGN KEY (idsocio) REFERENCES socio(idsocio),  
  FOREIGN KEY (idlibro) REFERENCES libro(idlibro)  
);



-- Procedure.

DELIMITER $$

CREATE PROCEDURE realizarPrestamo(
    IN p_loanID INT,          
    IN p_usuarioID INT,
    IN p_libroID INT,
    IN p_fechaPrestamo DATE,
    IN p_fechaDevolucion DATE
)
BEGIN
    DECLARE estadoLibro ENUM('disponible', 'prestado', 'reservado');
    
    SELECT estado INTO estadoLibro FROM libro WHERE idlibro = p_libroID;
    
    IF estadoLibro != 'disponible' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El libro no está disponible para préstamo';
    END IF;

   INSERT INTO prestamo (idprestamo, idsocio, idlibro, fechaprestamo, fechadevolucion, estado)
    VALUES (p_loanID, p_usuarioID, p_libroID, p_fechaPrestamo, p_fechaDevolucion, 'activo');
    
  UPDATE libro SET estado = 'prestado' WHERE idlibro = p_libroID;
    
END $$

DELIMITER ;


-- Insertar editoriales
INSERT INTO editorial (ideditorial, nombre, direccion, paginaweb) VALUES
(20, 'Editorial D', 'Calle Ejemplo 101', 'www.editoriald.com'),
(21, 'Editorial E', 'Av. de la Innovación 102', 'www.editoriale.com'),
(22, 'Editorial F', 'Calle Tecnológica 103', 'www.editorialf.com'),
(23, 'Editorial G', 'Calle Creativa 104', 'www.editorialg.com');

-- Insertar autores
INSERT INTO autor (idautor, nombre, nacionalidad) VALUES
(20, 'Elena Martínez', 'España'),
(21, 'Roberto Díaz', 'Chile'),
(22, 'Lucía Fernández', 'Colombia'),
(23, 'José Hernández', 'Perú'),
(24, 'Carlos Pérez', 'México'),
(25, 'Ana Torres', 'Argentina'),
(26, 'David Sánchez', 'Uruguay'),
(27, 'Mariana Castro', 'Ecuador');

-- Insertar libros
INSERT INTO libro (idlibro, ideditorial, fechapublicacion, titulo, genero, estado) VALUES
(20, 20, '2022-01-15', 'La Ciencia del Futuro', 'Ciencia', 'disponible'),
(21, 21, '2020-07-20', 'La Travesía del Tiempo', 'Aventura', 'prestado'),
(22, 22, '2023-05-10', 'Misterios del Universo', 'Ficción', 'reservado'),
(23, 23, '2021-10-11', 'El Viaje Interdimensional', 'Ciencia Ficción', 'disponible'),
(24, 20, '2022-08-15', 'El Último Susurro', 'Suspenso', 'prestado'),
(25, 21, '2020-11-25', 'Secretos de la Naturaleza', 'Documental', 'reservado'),
(26, 22, '2021-04-13', 'Voces del Pasado', 'Historia', 'disponible'),
(27, 23, '2023-09-01', 'El Poder de la Imaginación', 'Fantasía', 'prestado'),
(28, 20, '2022-12-01', 'La Guerra del Olvido', 'Aventura', 'disponible'),
(29, 21, '2020-05-03', 'El Refugio Secreto', 'Misterio', 'reservado');

-- Insertar socios
INSERT INTO socio (idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, imagenperfil, rol) VALUES
(20, 'Luis González', 'Calle Larga 101', '100200300', 'luisg@gmail.com', '1992-02-20', 'normal', '2023-06-01', 'perfil20.jpg', 'usuario'),
(21, 'Sofía Ruiz', 'Av. Central 202', '200300400', 'sofiar@gmail.com', '1988-05-25', 'estudiante', '2023-07-10', 'perfil21.jpg', 'usuario'),
(22, 'Gabriel López', 'Calle Rápida 303', '300400500', 'gabriell@gmail.com', '1990-08-10', 'profesor', '2022-10-15', 'perfil22.jpg', 'administrador'),
(23, 'Isabella Pérez', 'Calle Futura 404', '400500600', 'isabellap@gmail.com', '1995-12-30', 'normal', '2023-08-05', 'perfil23.jpg', 'usuario'),
(24, 'Oscar Sánchez', 'Av. Siempre Viva 505', '500600700', 'oscars@gmail.com', '1993-07-15', 'estudiante', '2023-04-25', 'perfil24.jpg', 'usuario'),
(25, 'Raquel Martínez', 'Calle Tranquila 606', '600700800', 'raquelm@gmail.com', '1987-03-18', 'profesor', '2022-12-10', 'perfil25.jpg', 'administrador'),
(26, 'Andrés Díaz', 'Calle del Sol 707', '700800900', 'andresd@gmail.com', '1998-10-05', 'normal', '2023-01-12', 'perfil26.jpg', 'usuario'),
(27, 'Paola Fernández', 'Calle Lluvia 808', '800900100', 'paolaf@gmail.com', '1994-11-17', 'estudiante', '2023-03-20', 'perfil27.jpg', 'usuario'),
(28, 'Manuel García', 'Calle Libertad 909', '900100200', 'manuelg@gmail.com', '1991-06-22', 'profesor', '2022-09-05', 'perfil28.jpg', 'administrador'),
(29, 'Victoria Soto', 'Calle de los Pinos 1010', '100200300', 'victorias@gmail.com', '1992-01-11', 'normal', '2023-05-18', 'perfil29.jpg', 'usuario');

-- Insertar contraseñas de usuarios
INSERT INTO usuariopassword (idusuario, hash_contrasena) VALUES
(20, 'hash_contrasena_luis'),
(21, 'hash_contrasena_sofia'),
(22, 'hash_contrasena_gabriel'),
(23, 'hash_contrasena_isabella'),
(24, 'hash_contrasena_oscar'),
(25, 'hash_contrasena_raquel'),
(26, 'hash_contrasena_andres'),
(27, 'hash_contrasena_paola'),
(28, 'hash_contrasena_manuel'),
(29, 'hash_contrasena_victoria');


-- Insertar editoriales
INSERT INTO editorial (ideditorial, nombre, direccion, paginaweb) VALUES
(24, 'Editorial H', 'Calle de los Libros 108', 'www.editorialh.com'),
(25, 'Editorial I', 'Calle de la Creatividad 109', 'www.editoriali.com'),
(26, 'Editorial J', 'Calle Literaria 110', 'www.editorialj.com'),
(27, 'Editorial K', 'Calle Fantasía 111', 'www.editorialk.com');

-- Insertar autores
INSERT INTO autor (idautor, nombre, nacionalidad) VALUES
(28, 'Laura García', 'México'),
(29, 'Fernando Gómez', 'Chile'),
(30, 'Patricia López', 'Colombia'),
(31, 'Juan Martínez', 'Perú'),
(32, 'Marta Ruiz', 'España'),
(33, 'Ricardo González', 'Argentina'),
(34, 'Eva Sánchez', 'Ecuador'),
(35, 'Juanita Pérez', 'Uruguay');

-- Insertar libros
INSERT INTO libro (idlibro, ideditorial, fechapublicacion, titulo, genero, estado) VALUES
(30, 24, '2024-01-10', 'Exploradores del Tiempo', 'Ciencia Ficción', 'disponible'),
(31, 25, '2022-12-15', 'El Secreto de las Sombras', 'Misterio', 'prestado'),
(32, 26, '2023-02-22', 'El Último Refugio', 'Suspenso', 'reservado'),
(33, 27, '2024-04-01', 'Aventuras del Más Allá', 'Aventura', 'disponible'),
(34, 24, '2023-11-30', 'El Poder del Pensamiento', 'Filosofía', 'prestado'),
(35, 25, '2022-08-05', 'Historias de la Tierra Lejana', 'Fantasía', 'reservado'),
(36, 26, '2023-10-10', 'La Ciudad Subterránea', 'Ciencia Ficción', 'disponible'),
(37, 27, '2023-06-21', 'El Último Guerrero', 'Aventura', 'prestado'),
(38, 24, '2024-03-18', 'Misterios del Alma', 'Terror', 'disponible'),
(39, 25, '2023-05-28', 'Un Viaje al Pasado', 'Historia', 'reservado');

-- Insertar socios
INSERT INTO socio (idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, imagenperfil, rol) VALUES
(30, 'Natalia Torres', 'Calle del Sol 112', '011220330', 'natalia.t@gmail.com', '1990-05-14', 'normal', '2023-07-01', 'perfil30.jpg', 'usuario'),
(31, 'Carlos Fernández', 'Calle Verde 113', '022330440', 'carlosf@gmail.com', '1985-09-25', 'profesor', '2022-12-02', 'perfil31.jpg', 'administrador'),
(32, 'Luciana González', 'Calle del Mar 114', '033440550', 'lucianag@gmail.com', '1997-11-09', 'estudiante', '2023-09-10', 'perfil32.jpg', 'usuario'),
(33, 'Eduardo Pérez', 'Calle Nueva 115', '044550660', 'eduardop@gmail.com', '1994-01-12', 'normal', '2023-06-01', 'perfil33.jpg', 'usuario'),
(34, 'Marcela Díaz', 'Calle de la Luna 116', '055660770', 'marcelad@gmail.com', '1993-04-17', 'estudiante', '2023-08-05', 'perfil34.jpg', 'usuario'),
(35, 'Alfonso González', 'Calle del Viento 117', '066770880', 'alfonsog@gmail.com', '1991-02-22', 'profesor', '2023-07-20', 'perfil35.jpg', 'administrador'),
(36, 'Gabriela Castro', 'Calle Alta 118', '077880990', 'gabyc@gmail.com', '1995-10-30', 'normal', '2023-02-14', 'perfil36.jpg', 'usuario'),
(37, 'Javier Rodríguez', 'Calle de la Esperanza 119', '088990100', 'javierr@gmail.com', '1992-03-19', 'profesor', '2022-11-11', 'perfil37.jpg', 'administrador'),
(38, 'Susana López', 'Calle del Bosque 120', '099100210', 'susanal@gmail.com', '1989-12-03', 'estudiante', '2023-04-01', 'perfil38.jpg', 'usuario'),
(39, 'Ricardo Ramírez', 'Calle Clara 121', '100210320', 'ricardor@gmail.com', '1996-08-15', 'normal', '2023-07-08', 'perfil39.jpg', 'usuario');

-- Insertar contraseñas de usuarios
INSERT INTO usuariopassword (idusuario, hash_contrasena) VALUES
(30, 'hash_contrasena_natalia'),
(31, 'hash_contrasena_carlos'),
(32, 'hash_contrasena_luciana'),
(33, 'hash_contrasena_eduardo'),
(34, 'hash_contrasena_marcel'),
(35, 'hash_contrasena_alfonso'),
(36, 'hash_contrasena_gabriela'),
(37, 'hash_contrasena_javier'),
(38, 'hash_contrasena_susana'),
(39, 'hash_contrasena_ricardo');


-- Insertar préstamos
INSERT INTO prestamo (idprestamo, idsocio, idlibro, fechaprestamo, fechadevolucion, estado) VALUES
(40, 30, 30, '2024-04-01', '2024-04-15', 'activo'),
(41, 31, 31, '2023-12-15', '2024-01-15', 'completado'),
(42, 32, 32, '2024-05-22', '2024-06-22', 'activo'),
(43, 33, 33, '2024-06-01', '2024-06-15', 'completado'),
(44, 34, 34, '2024-03-18', '2024-04-18', 'activo'),
(45, 35, 35, '2023-11-30', '2023-12-30', 'completado'),
(46, 36, 36, '2024-04-10', '2024-04-25', 'activo'),
(47, 37, 37, '2023-06-21', '2023-07-21', 'completado'),
(48, 38, 38, '2024-03-18', '2024-04-01', 'activo'),
(49, 39, 39, '2024-01-10', '2024-01-24', 'completado');

-- Insertar multas
INSERT INTO multa (idmulta, idprestamo, saldopagar, fechamulta, estado) VALUES
(50, 40, 15.00, '2024-04-16', 'pendiente'),
(51, 41, 10.00, '2024-01-16', 'pagada'),
(52, 42, 5.00, '2024-06-23', 'pendiente'),
(53, 43, 20.00, '2024-06-16', 'pagada'),
(54, 44, 8.00, '2024-04-19', 'pendiente'),
(55, 45, 12.00, '2023-12-31', 'pagada'),
(56, 46, 7.50, '2024-04-26', 'pendiente'),
(57, 47, 18.00, '2023-07-22', 'pagada'),
(58, 48, 25.00, '2024-04-02', 'pendiente'),
(59, 49, 30.00, '2024-01-25', 'pagada');

-- Insertar relaciones libro-autor
INSERT INTO libro_autor (idlibro, idautor) VALUES
(30, 28),
(31, 29),
(32, 30),
(33, 31),
(34, 32),
(35, 33),
(36, 34),
(37, 35),
(38, 36),
(39, 37);

-- Insertar reservas
INSERT INTO reserva (idreserva, idsocio, idlibro, fechareserva, estado) VALUES
(60, 30, 30, '2024-04-01', 'activa'),
(61, 31, 31, '2023-12-15', 'completada'),
(62, 32, 32, '2024-05-22', 'activa'),
(63, 33, 33, '2024-06-01', 'cancelada'),
(64, 34, 34, '2024-03-18', 'completada'),
(65, 35, 35, '2023-11-30', 'activa'),
(66, 36, 36, '2024-04-10', 'completada'),
(67, 37, 37, '2023-06-21', 'activa'),
(68, 38, 38, '2024-03-18', 'cancelada'),
(69, 39, 39, '2024-01-10', 'completada');
