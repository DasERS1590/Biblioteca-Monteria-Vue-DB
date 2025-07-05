-- Procedimiento para convertir reserva a préstamo cuando el libro esté disponible
CREATE PROCEDURE convertirReservaAPrestamo(
    IN p_reservaID INT,
    IN p_fechaPrestamo DATE,
    IN p_fechaDevolucion DATE
)
BEGIN
    DECLARE v_libroID INT;
    DECLARE v_socioID INT;
    DECLARE v_estadoLibro ENUM('disponible', 'prestado', 'reservado');
    DECLARE v_estadoReserva ENUM('activa', 'cancelada', 'completada');
    
    -- Obtener información de la reserva
    SELECT idlibro, idsocio, estado 
    INTO v_libroID, v_socioID, v_estadoReserva
    FROM reserva 
    WHERE idreserva = p_reservaID;
    
    -- Verificar que la reserva esté activa
    IF v_estadoReserva != 'activa' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La reserva no está activa';
    END IF;
    
    -- Obtener estado actual del libro
    SELECT estado INTO v_estadoLibro 
    FROM libro 
    WHERE idlibro = v_libroID;
    
    -- Verificar que el libro esté disponible
    IF v_estadoLibro != 'disponible' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El libro no está disponible para préstamo';
    END IF;
    
    -- Crear el préstamo
    INSERT INTO prestamo (idsocio, idlibro, fechaprestamo, fechadevolucion, estado)
    VALUES (v_socioID, v_libroID, p_fechaPrestamo, p_fechaDevolucion, 'activo');
    
    -- Actualizar estado del libro
    UPDATE libro SET estado = 'prestado' WHERE idlibro = v_libroID;
    
    -- Marcar reserva como completada
    UPDATE reserva SET estado = 'completada' WHERE idreserva = p_reservaID;
    
END; 