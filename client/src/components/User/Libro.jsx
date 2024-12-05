import React, { useState } from "react";

const Libro = () => {
  const [genre, setGenre] = useState("");
  const [author, setAuthor] = useState("");
  const [titulo, setTitulo] = useState("");
  const [books, setBooks] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    setLoading(true);
    setError(null);
    setBooks([]);

    // Construir la URL con los parámetros de búsqueda
    const queryParams = new URLSearchParams();
    if (genre) queryParams.append("genero", genre);
    if (author) queryParams.append("autor", author);
    if (titulo) queryParams.append("titulo", titulo);

    try {
      const response = await fetch(
        `http://localhost:4000/api/books?${queryParams.toString()}`
      );

      if (response.ok) {
        const data = await response.json();
        setBooks(data);
      } else {
        const errorData = await response.json();
        setError(errorData.message || "Error al buscar libros");
      }
    } catch (err) {
      setError("Error al conectar con el servidor");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: "20px", fontFamily: "Arial, sans-serif" }}>
      <h1>Buscar Libros</h1>
      <div style={{ marginBottom: "20px" }}>
        <label>
          Género:{" "}
          <input
            type="text"
            value={genre}
            onChange={(e) => setGenre(e.target.value)}
            placeholder="Ej: Ciencia Ficción"
          />
        </label>
        <br />
        <label>
          Autor:{" "}
          <input
            type="text"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            placeholder="Ej: Isaac Asimov"
          />
        </label>
        <br />
        <label>
          Título:{" "}
          <input
            type="text"
            value={titulo}
            onChange={(e) => setTitulo(e.target.value)}
            placeholder="Ej: Fundación"
          />
        </label>
      </div>
      <button onClick={handleSearch} disabled={loading}>
        {loading ? "Buscando..." : "Buscar"}
      </button>

      {error && <p style={{ color: "red" }}>{error}</p>}

      {books.length > 0 && (
        <div style={{ marginTop: "20px" }}>
          <h2>Resultados</h2>
          <table border="1" cellPadding="10" style={{ width: "100%", textAlign: "left" }}>
            <thead>
              <tr>
                <th>ID</th>
                <th>Título</th>
                <th>Género</th>
                <th>Estado</th>
                <th>Autor</th>
              </tr>
            </thead>
            <tbody>
              {books.map((book) => (
                <tr key={book.id_libro}>
                  <td>{book.id_libro}</td>
                  <td>{book.titulo}</td>
                  <td>{book.genero}</td>
                  <td>{book.estado}</td>
                  <td>{book.autor}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default Libro;