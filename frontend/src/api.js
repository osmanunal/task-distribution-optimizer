import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000/api';

// API istekleri için Axios örneğini yapılandır
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// API servis fonksiyonları
const taskService = {
  // Görev dağılımı planını getir
  getTaskDistribution: async () => {
    try {
      const response = await apiClient.get('/tasks/plan');
      return response.data;
    } catch (error) {
      console.error('Görev dağılım planı alınırken hata oluştu:', error);
      throw error;
    }
  },
};

export default taskService; 