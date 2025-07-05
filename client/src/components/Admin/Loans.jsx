import { useApiGet, useFilters } from "../../hooks/useApi";
import { loanService } from "../../services/loanService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import "../../styles/common/Admin.css";

const Loans = () => {
  const { filters, updateFilter, clearFilters } = useFilters({
    startdate: "",
    enddate: ""
  });

  const { data: loansData, loading, error, refetch } = useApiGet(
    () => (filters.startdate && filters.enddate)
      ? loanService.getActiveLoans(filters)
      : Promise.resolve({ loans: [] }),
    [filters.startdate, filters.enddate]
  );

  // Función para manejar errores de permisos
  const handlePermissionError = (error) => {
    if (!error) return null;
    
    if (error?.message?.includes("permissions") || error?.message?.includes("access")) {
      return "No tienes permisos para acceder a esta funcionalidad. Contacta al administrador.";
    }
    return error?.message || "Error al cargar los datos";
  };

  const filterConfig = [
    {
      name: "startdate",
      type: "date",
      placeholder: "Fecha de inicio",
      value: filters.startdate
    },
    {
      name: "enddate",
      type: "date",
      placeholder: "Fecha de fin",
      value: filters.enddate
    }
  ];

  const columns = [
    { key: 'idprestamo', label: 'ID Préstamo' },
    { key: 'idsocio', label: 'ID Socio' },
    { key: 'nombre_socio', label: 'Nombre del Socio' },
    { key: 'titulo_libro', label: 'Título del Libro' },
    { key: 'fechaprestamo', label: 'Fecha Préstamo' },
    { key: 'fechadevolucion', label: 'Fecha Devolución' },
    { key: 'estado', label: 'Estado', render: (value) => (
      <span className={`status ${value}`}>{value}</span>
    )}
  ];

  // Asegurar que loansArray sea siempre un array
  const loansArray = loansData?.loans ? (Array.isArray(loansData.loans) ? loansData.loans : [loansData.loans]) : [];

  return (
    <div className="admin-dashboard">
      <h1>Préstamos Activos</h1>
      <p className="description">Consulta todos los préstamos activos en un rango de fechas específico.</p>

      <FilterForm
        filters={filterConfig}
        onFilterChange={updateFilter}
        onSearch={refetch}
        onClear={clearFilters}
        loading={loading}
      />

      <DataTable
        columns={columns}
        data={loansArray}
        loading={loading}
        error={handlePermissionError(error)}
        emptyMessage="No se encontraron préstamos activos en el rango de fechas seleccionado."
        className="loans-table"
      />
    </div>
  );
};

export default Loans;
