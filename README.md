# Task Distribution Optimizer

## Proje Yapısı

```
.
├── frontend/                
│   ├── src/                
│   ├── public/             
│   └── package.json        
│
├── backend/                
│   ├── internal/           
│   ├── pkg/               
│   ├── server/            
│   ├── sync/              
│   ├── config/            
│   ├── migration/         
│   └── scripts/           
│
├── docker-compose.dev.yml 
└── Makefile              
```


## Kurulum

### Geliştirme Ortamı

1. Gereksinimleri yükleyin:
   - Docker ve Docker Compose
   - Node.js ve NPM
   - Go 1.22+
   - Make

2. Geliştirme ortamını başlatın:
```bash
make devdb  # PostgreSQL veritabanını Docker ile başlatır
```

### Make Komutları

#### Veritabanı Migrasyon Komutları
```bash
make migrate-up      # Migrasyonları uygula
```
```bash
make migrate-down    # Son migrasyonu geri al
```
```bash
make migrate-status  # Migrasyon durumunu göster
```
```bash
make migrate-create  # Yeni migrasyon dosyası oluştur
```

#### Servis Komutları
```bash
make add-emp       # Dummy çalışanları ekle
make sync-start    # providerdan verileri çek
make plan          # Görev planlama işlemini başlat
```

#### Frontend Komutları
```bash
cd frontend
npm install
npm start
```



## API Endpointleri

Backend servisi şu endpointleri sunar:

```
GET    /api/tasks/plan    # Görev planlama ve dağıtım optimizasyonu
```

