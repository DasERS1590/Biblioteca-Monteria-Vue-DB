import React, { useEffect } from "react";
import { useApiGet, useFilters } from "../../hooks/useApi";
import { fineService } from "../../services/fineService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import "../../styles/admin/Fines.css";

const Fines = () => {
  const { filters, updateFilter, clearFilters } = useFilters({
    busqueda: "",
    estado: "pendiente"
  });

  const { data: finesData, loading, error, refetch } = useApiGet(
    () => {
      if (filters.busqueda) {
        return fineService.searchFinesByUser({ busqueda: filters.busqueda });
      } else {
        return fineService.getPendingFines();
      }
    },
    [filters.busqueda, filters.estado]
  );

  // Asegurar que fines sea siempre un array
  const fines = finesData?.fines ? (Array.isArray(finesData.fines) ? finesData.fines : [finesData.fines]) : [];

  const handleSearch = () => {
    refetch();
  };

  const filterConfig = [
    {
      name: "busqueda",
      type: "text",
      placeholder: "Buscar por nombre o email",
      value: filters.busqueda
    },
    {
      name: "estado",
      type: "select",
      placeholder: "Estado",
      value: filters.estado,
      options: [
        { value: "pendiente", label: "Pendiente" },
        { value: "pagada", label: "Pagada" }
      ]
    }
  ];

  const columns = [
    { key: 'idmulta', label: 'ID Multa' },
    { key: 'idsocio', label: 'ID Socio' },
    { key: 'nombre_socio', label: 'Nombre del Socio' },
    { key: 'idprestamo', label: 'ID Préstamo' },
    { key: 'saldopagar', label: 'Saldo a Pagar', render: (value) => `$${value}` },
    { key: 'fechamulta', label: 'Fecha Multa' },
    { key: 'estado', label: 'Estado', render: (value) => (
      <span className={`status ${value}`}>{value}</span>
    )}
  ];

  return (
    <div className="admin-dashboard">
      <h1>Multas</h1>
      <p className="description">Consulta y filtra las multas de usuarios. Puedes buscar por nombre o email del usuario.</p>
      
      <FilterForm
        filters={filterConfig}
        onFilterChange={updateFilter}
        onSearch={handleSearch}
        onClear={clearFilters}
        loading={loading}
      />

      <DataTable
        columns={columns}
        data={fines}
        loading={loading}
        error={error}
        emptyMessage="No se encontraron multas."
        className="fines-table"
      />
    </div>
  );
};

export default Fines;
