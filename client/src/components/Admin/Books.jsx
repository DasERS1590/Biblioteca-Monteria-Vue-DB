import React, { useState, useEffect } from "react";
import { useApiGet, useApiMutation, useForm, useFilters } from "../../hooks/useApi";
import { bookService } from "../../services/bookService";
import { authorService } from "../../services/authorService";
import { editorialService } from "../../services/editorialService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import Modal from "../common/Modal";
import ConfirmDialog from "../common/ConfirmDialog";
import "../../styles/admin/Books.css";
import "../../styles/common/Modal.css";

function Books() {
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [selectedBook, setSelectedBook] = useState(null);
  const [authors, setAuthors] = useState([]);
  const [editorials, setEditorials] = useState([]);
  const [editBookDraft, setEditBookDraft] = useState(null);
  const [editFormReady, setEditFormReady] = useState(false);

  const { filters, updateFilter, clearFilters } = useFilters({
    titulo: "",
    genero: "",
    autor: "",
    editorial: ""
  });

  // Estado separado para el formulario de editar
  const { formData: editFormData, handleChange: handleEditChange, resetForm: resetEditForm, setFormData: setEditFormData } = useForm({
    titulo: "",
    genero: "",
    autor_id: "",
    editorial_id: "",
    anio_publicacion: "",
    estado: "disponible"
  });

  const { data: booksData, loading, error, refetch } = useApiGet(
    () => bookService.getBooks(filters),
    [filters]
  );

  const { execute: updateBook, loading: updateLoading, error: updateError, success: updateSuccess, reset: resetUpdate } = useApiMutation(
    (data) => bookService.updateBook(selectedBook?.id_libro, data)
  );

  const { execute: deleteBook, loading: deleteLoading, error: deleteError, success: deleteSuccess, reset: resetDelete } = useApiMutation(
    (bookId) => bookService.deleteBook(bookId)
  );

  const { data: authorsData } = useApiGet(
    () => authorService.getAuthors(),
    []
  );

  const { data: editorialsData } = useApiGet(
    () => editorialService.getEditorials(),
    []
  );

  React.useEffect(() => {
    if (authorsData?.authors) {
      setAuthors(authorsData.authors);
    }
    if (editorialsData?.editorials) {
      setEditorials(editorialsData.editorials);
    }
  }, [authorsData, editorialsData]);

  useEffect(() => {
    if (
      showEditModal &&
      editBookDraft &&
      authors.length > 0 &&
      editorials.length > 0
    ) {
      setEditFormData({
        titulo: editBookDraft.titulo || "",
        genero: editBookDraft.genero || "",
        autor_id:
          editBookDraft.autores && editBookDraft.autores.length > 0
            ? String(editBookDraft.autores[0])
            : "",
        editorial_id: editBookDraft.editorial_id
          ? String(editBookDraft.editorial_id)
          : "",
        anio_publicacion: editBookDraft.fecha_publicacion
          ? editBookDraft.fecha_publicacion.split("-")[0]
          : "",
        estado: editBookDraft.estado || "disponible",
      });
      setEditBookDraft(null);
    }
  }, [showEditModal, editBookDraft, authors, editorials, setEditFormData]);

  const handleUpdateBook = async (e) => {
    e.preventDefault();
    try {
      const bookData = {
        titulo: editFormData.titulo,
        genero: editFormData.genero,
        fechapublicacion: editFormData.anio_publicacion ? `${editFormData.anio_publicacion}-01-01` : "",
        ideditorial: parseInt(editFormData.editorial_id),
        idautores: [parseInt(editFormData.autor_id)]
      };
      
      await updateBook(bookData);
      resetEditForm();
      setShowEditModal(false);
      setSelectedBook(null);
      refetch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const handleDeleteBook = async () => {
    if (!selectedBook) {return;}

    try {
      await deleteBook(selectedBook.id_libro);
      setShowDeleteDialog(false);
      setSelectedBook(null);
      refetch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const openEditModal = async (book) => {
    setSelectedBook(book);
    setEditFormReady(false);
    setShowEditModal(true);

    try {
      const response = await bookService.getBookForEdit(book.id_libro);
      const bookData = response.book;

      // Espera a que autores y editoriales estén cargados
      const waitForData = () => {
        if (authors.length > 0 && editorials.length > 0) {
          setEditFormData({
            titulo: bookData.titulo || "",
            genero: bookData.genero || "",
            autor_id: bookData.autores && bookData.autores.length > 0 ? String(bookData.autores[0]) : "",
            editorial_id: bookData.editorial_id ? String(bookData.editorial_id) : "",
            anio_publicacion: bookData.fecha_publicacion ? bookData.fecha_publicacion.split('-')[0] : "",
            estado: bookData.estado || "disponible"
          });
          setEditFormReady(true);
        } else {
          setTimeout(waitForData, 50);
        }
      };
      waitForData();
    } catch (error) {
      setEditFormData({
        titulo: book.titulo || "",
        genero: book.genero || "",
        autor_id: "",
        editorial_id: "",
        anio_publicacion: "",
        estado: book.estado || "disponible"
      });
      setEditFormReady(true);
    }
    resetUpdate();
  };

  const openDeleteDialog = (book) => {
    setSelectedBook(book);
    setShowDeleteDialog(true);
    resetDelete();
  };

  const closeEditModal = () => {
    setShowEditModal(false);
    setSelectedBook(null);
    resetEditForm();
    resetUpdate();
  };

  const filterConfig = [
    {
      name: "titulo",
      placeholder: "Título del libro",
      value: filters.titulo
    },
    {
      name: "genero",
      placeholder: "Género",
      value: filters.genero
    },
    {
      name: "autor",
      placeholder: "Autor",
      value: filters.autor
    },
    {
      name: "editorial",
      placeholder: "Editorial",
      value: filters.editorial
    }
  ];

  const columns = [
    { key: 'id_libro', label: 'ID' },
    { key: 'titulo', label: 'Título' },
    { key: 'genero', label: 'Género' },
    { key: 'autor', label: 'Autor' },
    { key: 'editorial', label: 'Editorial' },
    { key: 'anio_publicacion', label: 'Año' },
    { key: 'estado', label: 'Estado', render: (value) => (
      <span className={`status ${value}`}>{value}</span>
    )},
    { 
      key: 'actions', 
      label: 'Acciones',
      render: (_, book) => (
        <div className="book-actions">
          <button 
            className="btn btn-primary btn-sm"
            onClick={() => openEditModal(book)}
          >
            Editar
          </button>
          <button 
            className="btn btn-danger btn-sm"
            onClick={() => openDeleteDialog(book)}
          >
            Eliminar
          </button>
        </div>
      )
    }
  ];

  const books = booksData?.books || [];

  return (
    <div className="books-dashboard">
      <div className="books-header">
        <h1>Gestión de Libros</h1>
      </div>
      
      <p className="books-description">Administra el catálogo de libros de la biblioteca.</p>
      
      <FilterForm
        filters={filterConfig}
        onFilterChange={updateFilter}
        onSearch={() => refetch()}
        onClear={clearFilters}
        loading={loading}
      />
      
      <DataTable
        columns={columns}
        data={books}
        loading={loading}
        error={error}
        emptyMessage="No se encontraron libros."
        className="books-table"
      />



      {/* Modal para editar libro */}
      <Modal
        isOpen={showEditModal}
        onClose={closeEditModal}
        title="Editar Libro"
        className="book-modal"
      >
        {!editFormReady ? (
          <div style={{ textAlign: "center", padding: "2rem" }}>Cargando datos...</div>
        ) : (
          <form onSubmit={handleUpdateBook} className="book-form">
            <div className="form-row">
              <div className="form-group">
                <label>Título:</label>
                <input
                  type="text"
                  name="titulo"
                  value={editFormData.titulo}
                  onChange={handleEditChange}
                  required
                />
              </div>
              <div className="form-group">
                <label>Género:</label>
                <input
                  type="text"
                  name="genero"
                  value={editFormData.genero}
                  onChange={handleEditChange}
                  required
                />
              </div>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Autor:</label>
                <select
                  name="autor_id"
                  value={editFormData.autor_id || ""}
                  onChange={handleEditChange}
                  required
                >
                  {editFormData.autor_id === "" && <option value="">Seleccionar autor</option>}
                  {authors.map((author) => (
                    <option key={author.id_autor} value={String(author.id_autor)}>
                      {author.nombre}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>Editorial:</label>
                <select
                  name="editorial_id"
                  value={editFormData.editorial_id || ""}
                  onChange={handleEditChange}
                  required
                >
                  {editFormData.editorial_id === "" && <option value="">Seleccionar editorial</option>}
                  {editorials.map((editorial) => (
                    <option key={editorial.id_editorial} value={String(editorial.id_editorial)}>
                      {editorial.nombre}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Año de Publicación:</label>
                <input
                  type="number"
                  name="anio_publicacion"
                  value={editFormData.anio_publicacion}
                  onChange={handleEditChange}
                  min="1900"
                  max={new Date().getFullYear()}
                  required
                />
              </div>
              <div className="form-group">
                <label>Estado:</label>
                <select
                  name="estado"
                  value={editFormData.estado}
                  onChange={handleEditChange}
                  required
                >
                  <option value="disponible">Disponible</option>
                  <option value="prestado">Prestado</option>
                  <option value="mantenimiento">Mantenimiento</option>
                </select>
              </div>
            </div>

            <div className="form-actions">
              <button
                type="submit"
                className="btn btn-primary"
                disabled={updateLoading}
              >
                {updateLoading ? "Actualizando..." : "Actualizar Libro"}
              </button>
              <button
                type="button"
                className="btn btn-secondary"
                onClick={closeEditModal}
              >
                Cancelar
              </button>
            </div>

            {updateError && <div className="error-message">{updateError}</div>}
            {updateSuccess && <div className="success-message">Libro actualizado exitosamente</div>}
          </form>
        )}
      </Modal>

      {/* Diálogo de confirmación para eliminar */}
      <ConfirmDialog
        isOpen={showDeleteDialog}
        title="Eliminar Libro"
        message={`¿Estás seguro de que quieres eliminar el libro "${selectedBook?.titulo}"? Esta acción no se puede deshacer.`}
        onConfirm={handleDeleteBook}
        onCancel={() => {
          setShowDeleteDialog(false);
          setSelectedBook(null);
        }}
        confirmText={deleteLoading ? "Eliminando..." : "Eliminar"}
        cancelText="Cancelar"
        type="danger"
        disabled={deleteLoading}
      />

      {deleteError && (
        <div className="error-message">{deleteError}</div>
      )}
      {deleteSuccess && (
        <div className="success-message">Libro eliminado exitosamente</div>
      )}
    </div>
  );
}

export default Books;
