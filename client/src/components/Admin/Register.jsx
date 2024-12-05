import React, { useState } from "react";
import axios from "axios";

function Register() {
  const [formData, setFormData] = useState({
    nombre: "",
    direccion: "",
    telefono: "",
    correo: "",
    fechaNacimiento: "",
    tipoSocio: "",
    contrasena: "",
    rol: "usuario", // Rol por defecto
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Manejar el cambio de cada campo en el formulario
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  // Manejar el envío del formulario
  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // Realizar la solicitud POST
      const response = await axios.post("http://localhost:4000/api/register", formData, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      // Mostrar mensaje de éxito
      setSuccess(response.data.message);
      setError(""); // Limpiar el mensaje de error si es exitoso
    } catch (err) {
      // Mostrar mensaje de error
      if (err.response) {
        setError(err.response.data.error || "Error desconocido");
      } else {
        setError("No se pudo conectar al servidor");
      }
      setSuccess(""); // Limpiar el mensaje de éxito si hay error
    }
  };

  return (
    <div>
      <h2>Registro de Usuario</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {success && <p style={{ color: "green" }}>{success}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="nombre">Nombre</label>
          <input
            type="text"
            id="nombre"
            name="nombre"
            value={formData.nombre}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="direccion">Dirección</label>
          <input
            type="text"
            id="direccion"
            name="direccion"
            value={formData.direccion}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="telefono">Teléfono</label>
          <input
            type="text"
            id="telefono"
            name="telefono"
            value={formData.telefono}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="correo">Correo Electrónico</label>
          <input
            type="email"
            id="correo"
            name="correo"
            value={formData.correo}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="fechaNacimiento">Fecha de Nacimiento</label>
          <input
            type="date"
            id="fechaNacimiento"
            name="fechaNacimiento"
            value={formData.fechaNacimiento}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="tipoSocio">Tipo de Socio</label>
          <select
            id="tipoSocio"
            name="tipoSocio"
            value={formData.tipoSocio}
            onChange={handleChange}
            required
          >
            <option value="socio">Socio</option>
            <option value="socio_activo">Socio Activo</option>
          </select>
        </div>
        <div>
          <label htmlFor="contrasena">Contraseña</label>
          <input
            type="password"
            id="contrasena"
            name="contrasena"
            value={formData.contrasena}
            onChange={handleChange}
            required
          />
        </div>
        <div>
          <label htmlFor="rol">Rol</label>
          <select
            id="rol"
            name="rol"
            value={formData.rol}
            onChange={handleChange}
            required
          >
            <option value="usuario">Usuario</option>
            <option value="administrador">Administrador</option>
          </select>
        </div>
        <button type="submit">Registrar</button>
      </form>
    </div>
  );
}

export default Register;
