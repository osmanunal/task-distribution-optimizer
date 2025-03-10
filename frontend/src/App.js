import React, { useState, useEffect } from 'react';
import { Container, Box, Typography, Button, CircularProgress, Alert } from '@mui/material';
import EmployeeTaskTable from './components/EmployeeTaskTable';
import taskService from './api';

function App() {
  const [loading, setLoading] = useState(false);
  const [taskData, setTaskData] = useState(null);
  const [error, setError] = useState(null);

  // Görev dağılım verilerini getir
  const fetchTaskDistribution = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const data = await taskService.getTaskDistribution();
      setTaskData(data);
    } catch (err) {
      setError('Görev dağılım verileri alınırken bir hata oluştu. Lütfen daha sonra tekrar deneyin.');
      console.error('API hata:', err);
    } finally {
      setLoading(false);
    }
  };

  // Sayfa yüklendiğinde verileri otomatik olarak getir
  useEffect(() => {
    fetchTaskDistribution();
  }, []);

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 4 }}>
        <Typography variant="h3" component="h1" gutterBottom>
          Görev Dağılım Optimizasyonu
        </Typography>
        
        <Box sx={{ mb: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button 
            variant="contained" 
            onClick={fetchTaskDistribution}
            disabled={loading}
          >
            {loading ? 'Yükleniyor...' : 'Yenile'}
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
            <CircularProgress />
          </Box>
        ) : (
          <EmployeeTaskTable taskData={taskData} />
        )}
      </Box>
    </Container>
  );
}

export default App; 