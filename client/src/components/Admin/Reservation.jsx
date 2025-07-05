import { useApiGet, useFilters } from "../../hooks/useApi";
import { reservationService } from "../../services/reservationService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import "../../styles/common/Admin.css";

const Reservation = () => {
  const { filters, updateFilter, clearFilters } = useFilters({
    idsocio: "",
    idlibro: "",
    fechareserva: "",
    nombre_socio: ""
  });

  const { data: reservationsData, loading, error, refetch } = useApiGet(
    () => reservationService.getActiveReservations(filters),
    [filters.idsocio, filters.idlibro, filters.fechareserva, filters.nombre_socio]
  );

  // Asegurar que reservations sea siempre un array
  const reservations = reservationsData?.reservations ? (Array.isArray(reservationsData.reservations) ? reservationsData.reservations : [reservationsData.reservations]) : [];

  const handleSearch = () => {
    refetch();
  };

  const filterConfig = [
    {
      name: "idsocio",
      type: "text",
      placeholder: "ID del Socio",
      value: filters.idsocio
    },
    {
      name: "idlibro",
      type: "text",
      placeholder: "ID del Libro",
      value: filters.idlibro
    },
    {
      name: "fechareserva",
      type: "date",
      placeholder: "Fecha de reserva",
      value: filters.fechareserva
    },
    {
      name: "nombre_socio",
      type: "text",
      placeholder: "Nombre del Socio",
      value: filters.nombre_socio
    }
  ];

  const columns = [
    { key: 'idreserva', label: 'ID Reserva' },
    { key: 'idsocio', label: 'ID Socio' },
    { key: 'nombre_socio', label: 'Nombre Socio' },
    { key: 'idlibro', label: 'ID Libro' },
    { key: 'titulo_libro', label: 'Título Libro' },
    { key: 'fechareserva', label: 'Fecha Reserva' },
    { key: 'estado_reserva', label: 'Estado Reserva' },
    { key: 'genero_libro', label: 'Género Libro' },
    { key: 'editorial', label: 'Editorial' },
    { key: 'telefono_socio', label: 'Teléfono Socio' },
    { key: 'correo_socio', label: 'Correo Socio' },
    { key: 'tiposocio', label: 'Tipo Socio' },
    { key: 'fechanacimiento', label: 'Fecha Nacimiento' },
    { key: 'fecharegistro', label: 'Fecha Registro' }
  ];

  return (
    <div className="admin-dashboard">
      <h1>Reservas Activas</h1>
      <p className="description">Consulta las reservas activas según los filtros proporcionados.</p>

      <FilterForm
        filters={filterConfig}
        onFilterChange={updateFilter}
        onSearch={handleSearch}
        onClear={clearFilters}
        loading={loading}
      />

      <DataTable
        columns={columns}
        data={reservations}
        loading={loading}
        error={error}
        emptyMessage="No se encontraron reservas activas."
        className="reservation-table"
      />
    </div>
  );
};

export default Reservation;
