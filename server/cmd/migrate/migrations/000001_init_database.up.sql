CREATE TABLE editorial (  
  ideditorial INT AUTO_INCREMENT PRIMARY KEY,  
  nombre VARCHAR(255),  
  direccion VARCHAR(255),  
  paginaweb VARCHAR(100)  
);  
  
CREATE TABLE libro (  
  idlibro INT AUTO_INCREMENT  PRIMARY KEY,  
  ideditorial INT,  
  fechapublicacion DATE,  
  titulo VARCHAR(255),  
  genero VARCHAR(100),  
  estado ENUM('disponible', 'prestado', 'reservado'),  
  FOREIGN KEY (ideditorial) REFERENCES editorial(ideditorial) 
);  
  
CREATE TABLE autor (  
  idautor INT AUTO_INCREMENT PRIMARY KEY,  
  nombre VARCHAR(255),  
  nacionalidad VARCHAR(100)  
);  
  
CREATE TABLE socio (  
  idsocio INT AUTO_INCREMENT  PRIMARY KEY,  
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
  idprestamo INT AUTO_INCREMENT  PRIMARY KEY,  
  idsocio INT,  
  idlibro INT,  
  fechaprestamo DATE,  
  fechadevolucion DATE,  
  estado ENUM('activo', 'completado'),  
  FOREIGN KEY (idsocio) REFERENCES socio(idsocio),  
  FOREIGN KEY (idlibro) REFERENCES libro(idlibro)  
);  
  
CREATE TABLE  multa (  
  idmulta INT AUTO_INCREMENT  PRIMARY KEY,  
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
  idreserva INT AUTO_INCREMENT  PRIMARY KEY,  
  idsocio INT,  
  idlibro INT,  
  fechareserva DATE,  
  estado ENUM('activa', 'cancelada', 'completada'),  
  FOREIGN KEY (idsocio) REFERENCES socio(idsocio),  
  FOREIGN KEY (idlibro) REFERENCES libro(idlibro)  
);


