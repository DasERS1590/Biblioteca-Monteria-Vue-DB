import { useState, useEffect, useCallback } from 'react';

// Hook para peticiones GET
export const useApiGet = (apiCall, dependencies = []) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await apiCall();
      setData(result);
    } catch (err) {
      setError(err.message || 'Error al cargar los datos');
    } finally {
      setLoading(false);
    }
  }, [apiCall]);

  useEffect(() => {
    fetchData();
  }, dependencies);

  const refetch = useCallback(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    loading,
    error,
    refetch
  };
};

// Hook específico para préstamos que maneja 404 como datos vacíos
export const useLoansGet = (apiCall, dependencies = []) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await apiCall();
      setData(result);
    } catch (err) {
      // Si es un 404 con mensaje de "No hay préstamos", tratarlo como éxito con datos vacíos
      if (err.status === 404 && err.message && 
          (err.message.includes("No hay préstamos activos") || 
           err.message.includes("No hay préstamos completados") ||
           err.message.includes("No hay préstamos registrados"))) {

        setData({ loansactive: [], loanscomplete: [], loanshystory: [] });
        setError(null);
      } else {
        setError(err.message || 'Error al cargar los datos');
      }
    } finally {
      setLoading(false);
    }
  }, [apiCall]);

  useEffect(() => {
    fetchData();
  }, dependencies);

  const refetch = useCallback(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    loading,
    error,
    refetch
  };
};

// Hook para peticiones POST/PUT/DELETE
export const useApiMutation = (apiCall) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  const execute = useCallback(async (data) => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(false);
      
      const result = await apiCall(data);
      setSuccess(true);
      return result;
    } catch (err) {
      setError(err.message || 'Error en la operación');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiCall]);

  const reset = useCallback(() => {
    setError(null);
    setSuccess(false);
  }, []);

  return {
    execute,
    loading,
    error,
    success,
    reset
  };
};

// Hook para manejar formularios
export const useForm = (initialState = {}) => {
  const [formData, setFormData] = useState(initialState);

  const handleChange = useCallback((e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  }, []);

  const resetForm = useCallback(() => {
    setFormData(initialState);
  }, [initialState]);

  const setFormDataDirectly = useCallback((data) => {
    setFormData(data);
  }, []);

  return {
    formData,
    handleChange,
    resetForm,
    setFormData: setFormDataDirectly
  };
};

// Hook para manejar filtros
export const useFilters = (initialFilters = {}) => {
  const [filters, setFilters] = useState(initialFilters);

  const updateFilter = useCallback((nameOrEvent, value) => {
    // Si se pasa un evento (e.target), extraer name y value
    if (nameOrEvent && nameOrEvent.target) {
      const { name, value } = nameOrEvent.target;
      setFilters(prev => ({
        ...prev,
        [name]: value
      }));
    } else {
      // Si se pasan name y value directamente
      setFilters(prev => ({
        ...prev,
        [nameOrEvent]: value
      }));
    }
  }, []);

  const clearFilters = useCallback(() => {
    setFilters(initialFilters);
  }, [initialFilters]);

  const setFiltersDirectly = useCallback((newFilters) => {
    setFilters(newFilters);
  }, []);

  return {
    filters,
    updateFilter,
    clearFilters,
    setFilters: setFiltersDirectly
  };
}; 