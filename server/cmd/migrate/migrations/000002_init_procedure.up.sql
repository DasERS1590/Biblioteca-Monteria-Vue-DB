CREATE PROCEDURE realizarPrestamo(
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

    INSERT INTO prestamo (idsocio, idlibro, fechaprestamo, fechadevolucion, estado)
    VALUES (p_usuarioID, p_libroID, p_fechaPrestamo, p_fechaDevolucion, 'activo');
    
    UPDATE libro SET estado = 'prestado' WHERE idlibro = p_libroID;
    
END;

