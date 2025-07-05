import React from 'react';
import '../../styles/common/DataTable.css';

const DataTable = ({ 
  columns, 
  data, 
  loading, 
  error, 
  emptyMessage = "No hay datos disponibles.",
  className = "",
  onRowClick = null 
}) => {
  if (loading) {
    return (
      <div className={`data-table-container ${className}`}>
        <div className="loading-state">
          <div className="spinner"></div>
          <p>Cargando datos...</p>
        </div>
      </div>
    );
  }

  if (error) {
    // Verificar si es un error de permisos
    const isPermissionError = error.includes("permisos") || 
                             error.includes("permissions") || 
                             error.includes("access") ||
                             error.includes("403");
    
    return (
      <div className={`data-table-container ${className}`}>
        <div className="error-state">
          <div className="error-icon">⚠️</div>
          <h3>Error al cargar los datos</h3>
          <p className="error-message">{error}</p>
          {isPermissionError && (
            <div className="permission-help">
              <p><strong>Sugerencias:</strong></p>
              <ul>
                <li>Verifica que tengas los permisos necesarios para esta funcionalidad</li>
                <li>Si eres usuario, asegúrate de estar en la sección correcta</li>
                <li>Contacta al administrador si necesitas acceso adicional</li>
              </ul>
            </div>
          )}
        </div>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className={`data-table-container ${className}`}>
        <div className="empty-state">
          <div className="empty-icon">📋</div>
          <p>{emptyMessage}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`data-table-container ${className}`}>
      <div className="data-table">
        <table>
          <thead>
            <tr>
              {columns.map((column) => (
                <th key={column.key}>{column.label}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((row, index) => (
              <tr 
                key={row.id || index} 
                onClick={() => onRowClick && onRowClick(row)}
                className={onRowClick ? 'clickable-row' : ''}
              >
                {columns.map((column) => (
                  <td key={column.key}>
                    {column.render ? column.render(row[column.key], row) : row[column.key]}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DataTable; 