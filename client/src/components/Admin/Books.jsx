import React, { useState, useEffect } from "react";
import axios from "axios";

const Books = () => {
  const [books, setBooks] = useState([]); // Estado para guardar los libros
  const [filters, setFilters] = useState({ estado: "", editorial: "" }); // Estado para los filtros
  const [loading, setLoading] = useState(true); // Estado para el indicador de carga
  const [error, setError] = useState(null); // Estado para los errores

  // Función para obtener los libros desde la API
  const fetchBooks = async () => {
    try {
      setLoading(true); // Inicia la carga
      setError(null); // Limpia errores previos
      const response = await axios.get("http://localhost:4000/api/admin/books", {
        params: filters, // Agrega los filtros a la consulta
      });
      setBooks(response.data); // Guarda los datos en el estado
    } catch (err) {
      setError("Error al obtener los libros. Por favor, intenta más tarde.");
    } finally {
      setLoading(false); // Finaliza la carga
    }
  };

  // Ejecuta fetchBooks al cargar el componente o cuando los filtros cambien
  useEffect(() => {
    fetchBooks();
  }, [filters]);

  // Maneja cambios en los inputs de filtros
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFilters((prevFilters) => ({
      ...prevFilters,
      [name]: value,
    }));
  };

  return (
    <div style={styles.container}>
      <h1 style={styles.title}>Admin Dashboard - Libros</h1>

      {/* Filtros en la parte superior */}
      <div style={styles.filterContainer}>
        <label style={styles.label}>
          Estado:
          <input
            type="text"
            name="estado"
            value={filters.estado}
            onChange={handleInputChange}
            placeholder="Ej: disponible"
            style={styles.input}
          />
        </label>
        <label style={styles.label}>
          Editorial:
          <input
            type="text"
            name="editorial"
            value={filters.editorial}
            onChange={handleInputChange}
            placeholder="Ej: Santillana"
            style={styles.input}
          />
        </label>
        <button onClick={fetchBooks} style={styles.button}>
          Buscar
        </button>
      </div>

      {/* Indicador de carga */}
      {loading && <p style={styles.loading}>Cargando libros...</p>}

      {/* Error */}
      {error && <p style={styles.error}>{error}</p>}

      {/* Resultados debajo de los filtros */}
      {!loading && !error && (
        <div style={styles.bookList}>
          <h2 style={styles.sectionTitle}>Libros</h2>
          {books.length === 0 ? (
            <p style={styles.noResults}>No se encontraron libros con los filtros seleccionados.</p>
          ) : (
            <table style={styles.table}>
              <thead>
                <tr>
                  <th style={styles.th}>Título</th>
                  <th style={styles.th}>Género</th>
                  <th style={styles.th}>Estado</th>
                  <th style={styles.th}>Editorial</th>
                </tr>
              </thead>
              <tbody>
                {books.map((book) => (
                  <tr key={book.id} style={styles.row}>
                    <td style={styles.td}>{book.title}</td>
                    <td style={styles.td}>{book.genre}</td>
                    <td style={styles.td}>{book.status}</td>
                    <td style={styles.td}>{book.editorial}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  );
};

// Estilos embebidos
const styles = {
  container: {
    fontFamily: "Arial, sans-serif",
    padding: "20px",
    maxWidth: "1000px",
    margin: "0 auto",
    backgroundColor: "#f9f9f9",
    borderRadius: "8px",
    boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
  },
  title: {
    fontSize: "24px",
    color: "#333",
    marginBottom: "20px",
    textAlign: "center",
  },
  filterContainer: {
    display: "flex",
    gap: "15px",
    marginBottom: "20px",
    justifyContent: "center",
    alignItems: "center",
  },
  label: {
    display: "flex",
    flexDirection: "column",
    fontSize: "14px",
    color: "#333",
  },
  input: {
    padding: "10px",
    fontSize: "14px",
    borderRadius: "4px",
    border: "1px solid #ccc",
    marginTop: "5px",
  },
  button: {
    padding: "10px 20px",
    fontSize: "14px",
    backgroundColor: "#007bff",
    color: "#fff",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer",
    alignSelf: "flex-end",
  },
  buttonHover: {
    backgroundColor: "#0056b3",
  },
  loading: {
    fontSize: "16px",
    color: "#007bff",
    textAlign: "center",
  },
  error: {
    fontSize: "14px",
    color: "red",
    textAlign: "center",
  },
  bookList: {
    marginTop: "20px",
  },
  sectionTitle: {
    fontSize: "20px",
    color: "#333",
    marginBottom: "10px",
    textAlign: "center",
  },
  noResults: {
    fontSize: "14px",
    color: "#555",
    textAlign: "center",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
    marginTop: "20px",
  },
  th: {
    backgroundColor: "#007bff",
    color: "#fff",
    padding: "10px",
    textAlign: "left",
  },
  td: {
    border: "1px solid #ddd",
    padding: "10px",
    textAlign: "left",
  },
  row: {
    "&:nth-child(even)": {
      backgroundColor: "#f2f2f2",
    },
  },
};

export default Books;
