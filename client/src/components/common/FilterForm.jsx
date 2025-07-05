const FilterForm = ({ 
  filters, 
  onFilterChange, 
  onSearch, 
  onClear, 
  loading = false,
  className = "" 
}) => {
  const handleSubmit = (e) => {
    e.preventDefault();
    onSearch();
  };

  const renderField = (filter) => {
    if (filter.type === 'select') {
      return (
        <select
          key={filter.name}
          name={filter.name}
          value={filter.value || ''}
          onChange={(e) => onFilterChange(filter.name, e.target.value)}
          className="input"
        >
          <option value="">{filter.placeholder}</option>
          {filter.options?.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      );
    }

    return (
      <input
        key={filter.name}
        type={filter.type || 'text'}
        name={filter.name}
        value={filter.value || ''}
        onChange={(e) => onFilterChange(filter.name, e.target.value)}
        placeholder={filter.placeholder}
        className="input"
      />
    );
  };

  return (
    <form onSubmit={handleSubmit} className={`filters ${className}`}>
      {filters.map(renderField)}
      <button 
        type="submit" 
        className="btn" 
        disabled={loading}
      >
        {loading ? 'Buscando...' : 'Buscar'}
      </button>
      {onClear && (
        <button 
          type="button" 
          className="btn btn-secondary" 
          onClick={onClear}
          disabled={loading}
        >
          Limpiar
        </button>
      )}
    </form>
  );
};

export default FilterForm; 