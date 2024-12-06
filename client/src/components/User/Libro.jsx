import React, { useState } from "react";

const Libro = () => {
  const [genre, setGenre] = useState("");
  const [author, setAuthor] = useState("");
  const [titulo, setTitulo] = useState("");
  const [books, setBooks] = useState([]);
  const [selectedBook, setSelectedBook] = useState(null); // Libro seleccionado para préstamo
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [loanDates, setLoanDates] = useState({
    fechaPrestamo: "",
    fechaDevolucion: "",
  });
  const [loanMessage, setLoanMessage] = useState(""); // Mensaje de préstamo
  const [showModal, setShowModal] = useState(false); // Controla la visibilidad del modal

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

  // Manejar la selección de un libro
  const handleBookSelect = (book) => {
    setSelectedBook(book);
    setShowModal(true); // Mostrar el modal al seleccionar un libro
  };

  // Manejar el envío del formulario de préstamo
  const handleLoanSubmit = async (e) => {
    e.preventDefault();

   const user = JSON.parse(localStorage.getItem("user"));

    if (!user) {
      setError("Usuario no encontrado. Asegúrese de estar registrado.");
      return;
    }

    const loanData = {
      usuario_id: user.id,
      libro_id: selectedBook.id_libro,
      fecha_prestamo: loanDates.fechaPrestamo,
      fecha_devolucion: loanDates.fechaDevolucion,
    };

    try {
      const response = await fetch("http://localhost:4000/api/loans", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(loanData),
      });

      if (response.ok) {
        setLoanMessage("Préstamo realizado con éxito");
      } else {
        const errorData = await response.json();
        setLoanMessage(errorData.message || "Error al realizar el préstamo");
      }
    } catch (err) {
      setLoanMessage("Error al conectar con el servidor");
    }
  };

  // Cerrar el modal
  const closeModal = () => {
    setShowModal(false);
    setLoanMessage(""); // Limpiar mensaje al cerrar modal
    setLoanDates({ fechaPrestamo: "", fechaDevolucion: "" }); // Limpiar fechas
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
                <th>Acción</th>
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
                  <td>
                    <button onClick={() => handleBookSelect(book)}>Seleccionar</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showModal && selectedBook && (
        <div
          style={{
            position: "fixed",
            top: "0",
            left: "0",
            width: "100%",
            height: "100%",
            backgroundColor: "rgba(0,0,0,0.5)",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <div
            style={{
              backgroundColor: "white",
              padding: "20px",
              borderRadius: "5px",
              minWidth: "300px",
              textAlign: "center",
            }}
          >
            <h2>Formulario de Préstamo</h2>
            <p><strong>Libro Seleccionado:</strong> {selectedBook.titulo}</p>
            <form onSubmit={handleLoanSubmit}>
              <label>
                Fecha de Préstamo:
                <input
                  type="date"
                  value={loanDates.fechaPrestamo}
                  onChange={(e) => setLoanDates({ ...loanDates, fechaPrestamo: e.target.value })}
                  required
                />
              </label>
              <br />
              <label>
                Fecha de Devolución:
                <input
                  type="date"
                  value={loanDates.fechaDevolucion}
                  onChange={(e) => setLoanDates({ ...loanDates, fechaDevolucion: e.target.value })}
                  required
                />
              </label>
              <br />
              <button type="submit">Realizar Préstamo</button>
            </form>
            {loanMessage && (
              <p style={{ marginTop: "20px", color: loanMessage.includes("Error") ? "red" : "green" }}>
                {loanMessage}
              </p>
            )}
            <button onClick={closeModal} style={{ marginTop: "10px" }}>Cerrar</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Libro;
