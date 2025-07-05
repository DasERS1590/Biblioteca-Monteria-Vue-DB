
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS users_permissions (
    user_id INT NOT NULL,              
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id)
        REFERENCES socio(idsocio)
        ON DELETE CASCADE,
    FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);

-- 3) Carga inicial de permisos
INSERT INTO permissions (code)
VALUES 
    ('books:read'),
    ('books:write'),
    ('books:delete'),
    ('loans:create'),
    ('loans:view'),
    ('loans:manage'),
    ('users:view'),
    ('users:manage'),
    ('reports:generate'),
    ('fines:read'),
    ('fines:create'),
    ('publishers:read'), 
    ('publishers:create'), 
    ('authors:read'), 
    ('authors:create'), 
    ('reservations:create'),
    ('reservations:view');

