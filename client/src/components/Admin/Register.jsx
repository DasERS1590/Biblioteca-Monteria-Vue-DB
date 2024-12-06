import React, { useState } from 'react';
import axios from 'axios';

function Register() {
  const [formData, setFormData] = useState({
    nombre: '',
    direccion: '',
    telefono: '',
    correo: '',
    fecha_nacimiento: '', // Cambiado a fecha_nacimiento
    tipo_socio: 'normal', // Cambiado a tipo_socio
    contrasena: '',
    rol: 'usuario',
  });

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Maneja los cambios de los campos del formulario
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  // Maneja el envío del formulario
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    // Validar que todos los campos estén completos
    if (
      !formData.nombre ||
      !formData.direccion ||
      !formData.telefono ||
      !formData.correo ||
      !formData.fecha_nacimiento ||
      !formData.contrasena
    ) {
      setError('Todos los campos son obligatorios');
      return;
    }

    try {
      // Realizar solicitud POST al backend con la URL correcta
      const response = await axios.post('http://localhost:4000/api/register', formData, {
        headers: {
          'Content-Type': 'application/json',
        },
      });

      // Respuesta exitosa
      setSuccess(`Usuario registrado exitosamente. ID: ${response.data.id}`);
    } catch (error) {
      // Manejo de errores
      if (error.response) {
        setError(error.response.data.message || 'Error al registrar el usuario');
      } else {
        setError('Error en la conexión con el servidor');
      }
    }
  };

  return (
    <div>
      <h2>Registro de Usuario</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Nombre:</label>
          <input
            type="text"
            name="nombre"
            value={formData.nombre}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Dirección:</label>
          <input
            type="text"
            name="direccion"
            value={formData.direccion}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Teléfono:</label>
          <input
            type="text"
            name="telefono"
            value={formData.telefono}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Correo:</label>
          <input
            type="email"
            name="correo"
            value={formData.correo}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Fecha de Nacimiento:</label>
          <input
            type="date"
            name="fecha_nacimiento" // Cambiado a fecha_nacimiento
            value={formData.fecha_nacimiento}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Tipo de Socio:</label>
          <select
            name="tipo_socio" // Cambiado a tipo_socio
            value={formData.tipo_socio}
            onChange={handleChange}
          >
            <option value="normal">Normal</option>
            <option value="estudiante">Estudiante</option>
            <option value="profesor">Profesor</option>
          </select>
        </div>
        <div>
          <label>Contraseña:</label>
          <input
            type="password"
            name="contrasena"
            value={formData.contrasena}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Rol:</label>
          <select
            name="rol"
            value={formData.rol}
            onChange={handleChange}
          >
            <option value="usuario">Usuario</option>
            <option value="administrador">Administrador</option>
          </select>
        </div>
        <button type="submit">Registrar </button>
      </form>

      {error && <div style={{ color: 'red' }}>{error}</div>}
      {success && <div style={{ color: 'green' }}>{success}</div>}
    </div>
  );
}

export default Register;