import React, { useState } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { 
  Box, 
  Typography, 
  Tabs, 
  Tab, 
  Paper, 
  Accordion, 
  AccordionSummary, 
  AccordionDetails,
  Chip
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

// Çalışan görev dağılımı tablosu
const EmployeeTaskTable = ({ taskData }) => {
  const [tabIndex, setTabIndex] = useState(0);

  // Veri yoksa boş içerik göster
  if (!taskData || !taskData.workloads || taskData.workloads.length === 0) {
    return (
      <Box sx={{ p: 2 }}>
        <Typography variant="h6">Görev dağılımı verisi bulunamadı</Typography>
      </Box>
    );
  }

  // Sekmeler arasında geçiş için handler
  const handleTabChange = (event, newValue) => {
    setTabIndex(newValue);
  };

  // Görev tablosu sütunları
  const taskColumns = [
    { field: 'task_id', headerName: 'Görev ID', width: 100 },
    { field: 'task_name', headerName: 'Görev Adı', width: 300 },
    { field: 'task_external_id', headerName: 'Dış ID', width: 100 },
    { field: 'duration', headerName: 'Süre (Saat)', width: 120 },
  ];

  // Haftalık plan tablosu sütunları
  const weeklyColumns = [
    { field: 'week_number', headerName: 'Hafta No', width: 100 },
    { field: 'hours', headerName: 'Çalışma Saati', width: 150 },
  ];

  return (
    <Box sx={{ width: '100%' }}>
      <Typography variant="h4" sx={{ mb: 2 }}>
        Görev Dağılım Planı - Toplam {taskData.total_weeks} Hafta
      </Typography>
      
      <Box sx={{ mb: 4 }}>
        <Tabs value={tabIndex} onChange={handleTabChange}>
          <Tab label="Çalışan Bazlı Görünüm" />
          <Tab label="Haftalık İş Yükü" />
        </Tabs>
      </Box>

      {tabIndex === 0 ? (
        // Çalışan bazlı görünüm
        <Box>
          {taskData.workloads.map((employee) => (
            <Accordion key={employee.employee_id} sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography variant="h6">
                  {employee.employee_name} 
                  <Chip 
                    label={`${employee.total_hours} saat`}
                    color="primary"
                    size="small"
                    sx={{ ml: 2 }}
                  />
                  <Chip 
                    label={`Zorluk: ${employee.difficulty}`}
                    color="secondary"
                    size="small"
                    sx={{ ml: 1 }}
                  />
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                <Typography variant="subtitle1" sx={{ mb: 2 }}>
                  Atanan Görevler
                </Typography>
                <Box sx={{ height: 400, width: '100%' }}>
                  <DataGrid
                    rows={employee.assignments.map((task, index) => ({
                      id: index,
                      ...task
                    }))}
                    columns={taskColumns}
                    pageSize={5}
                    rowsPerPageOptions={[5, 10]}
                    disableSelectionOnClick
                  />
                </Box>
                
                <Typography variant="subtitle1" sx={{ mb: 2, mt: 4 }}>
                  Haftalık Plan
                </Typography>
                <Box sx={{ height: 300, width: '100%' }}>
                  <DataGrid
                    rows={employee.weekly_plan.map((week, index) => ({
                      id: index,
                      ...week
                    }))}
                    columns={weeklyColumns}
                    pageSize={5}
                    rowsPerPageOptions={[5, 10]}
                    disableSelectionOnClick
                  />
                </Box>
              </AccordionDetails>
            </Accordion>
          ))}
        </Box>
      ) : (
        // Haftalık görünüm
        <Paper elevation={3} sx={{ p: 2 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Haftalık Genel Bakış
          </Typography>
          {/* Burada haftalık görünüm için başka bir tablo eklenebilir */}
          {/* Örneğin hafta bazında toplam iş yükü dağılımı */}
        </Paper>
      )}
    </Box>
  );
};

export default EmployeeTaskTable; 