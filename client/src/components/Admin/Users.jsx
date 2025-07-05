import { useApiGet, useFilters } from "../../hooks/useApi";
import { userService } from "../../services/userService";
import DataTable from "../common/DataTable";
import "../../styles/common/Admin.css";

const User = () => {
  const { filters, updateFilter, clearFilters } = useFilters({
    tiposocio: ""
  });

  // Usa useApiGet para refrescar automáticamente al cambiar el filtro
  const { data: usersData, loading, error, refetch } = useApiGet(
    () => filters.tiposocio ? userService.getUsersByType(filters) : Promise.resolve({ usertype: [] }),
    [filters.tiposocio]
  );

  const handleUserTypeChange = (e) => {
    updateFilter("tiposocio", e.target.value);
  };

  const columns = [
    { key: 'id', label: 'ID' },
    { key: 'nombre', label: 'Nombre' },
    { key: 'direccion', label: 'Dirección' },
    { key: 'telefono', label: 'Teléfono' },
    { key: 'correo', label: 'Correo' },
    { key: 'tiposocio', label: 'Tipo de Socio' }
  ];

  const usersArray = Array.isArray(usersData?.usertype) ? usersData.usertype : [];

  return (
    <div className="admin-dashboard">
      <h1>Usuarios</h1>
      <p className="description">Selecciona un tipo de socio para ver los usuarios correspondientes.</p>
      
      <div className="filters">
        <select 
          value={filters.tiposocio} 
          onChange={handleUserTypeChange} 
          className="input"
        >
          <option value="">Seleccione un tipo de socio</option>
          <option value="normal">Normal</option>
          <option value="estudiante">Estudiante</option>
          <option value="profesor">Profesor</option>
        </select>
        <button 
          onClick={refetch} 
          className="btn"
          disabled={loading || !filters.tiposocio}
        >
          {loading ? 'Buscando...' : 'Buscar'}
        </button>
        <button 
          onClick={clearFilters} 
          className="btn btn-secondary"
          disabled={loading}
        >
          Limpiar
        </button>
      </div>

      <DataTable
        columns={columns}
        data={usersArray}
        loading={loading}
        error={error}
        emptyMessage="No se encontraron usuarios."
        className="users-table"
      />
    </div>
  );
};

export default User;