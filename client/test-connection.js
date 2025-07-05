// Script para probar la conexión al backend
const axios = require('axios');

const API_BASE_URL = 'http://localhost:4000/v1/api';

async function testConnection() {
  console.log('🔍 Probando conexión al backend...');
  console.log(`URL: ${API_BASE_URL}`);
  
  try {
    // Probar el healthcheck primero
    const healthResponse = await axios.get('http://localhost:4000/v1/healthcheck');
    console.log('✅ Healthcheck exitoso:', healthResponse.data);
    
    // Probar el endpoint de libros (sin autenticación para ver si responde)
    try {
      const booksResponse = await axios.get(`${API_BASE_URL}/books`);
      console.log('✅ Endpoint de libros responde:', booksResponse.status);
    } catch (booksError) {
      if (booksError.response?.status === 401) {
        console.log('✅ Endpoint de libros funciona (requiere autenticación)');
      } else {
        console.log('❌ Error en endpoint de libros:', booksError.message);
      }
    }
    
  } catch (error) {
    console.log('❌ Error de conexión:', error.message);
    console.log('💡 Asegúrate de que el servidor esté corriendo en el puerto 4000');
  }
}

testConnection(); 