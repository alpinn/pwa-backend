# PWA Backend API

Backend API untuk Progressive Web Application yang dibangun dengan Go, PostgreSQL, dan Swagger documentation.

## Prerequisites

- **Go** 1.21 atau lebih tinggi
- **PostgreSQL** 13 atau lebih tinggi (atau gunakan Docker Compose)
- **Docker** dan **Docker Compose** (opsional, untuk menjalankan dengan containerization)
- **Git** untuk version control

## Setup Environment

### 1. Clone Repository
```bash
git clone <repository-url>
cd pwa-backend
```

### 2. Setup Environment Variables
Copy file `.env.example` ke `.env` dan sesuaikan nilai-nilainya:

```bash
cp .env.example .env
```

Edit `.env` dengan konfigurasi yang sesuai:
```env
DATABASE_URL=postgresql://username:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET=your-secret-key-here
PORT=8080
```

### 3. Install Dependencies
```bash
go mod download
go mod tidy
```

## Cara Menjalankan

### Option 1: Menggunakan Docker Compose (Recommended)

Cara paling mudah untuk menjalankan project beserta database:

**Step 1: Build Docker Image**
```bash
docker-compose build
```

**Step 2: Start Services**
```bash
docker-compose up -d
```

Ini akan:
- Membuat PostgreSQL container
- Build dan run aplikasi backend
- Expose API di `http://localhost:8080`

Untuk stop:
```bash
docker-compose down
```

### Option 2: Menjalankan Secara Lokal

**Step 1: Setup Database**

Pastikan PostgreSQL berjalan, kemudian buat database:
```bash
createdb pwa_db
```

**Step 2: Run Migrations**

Jalankan SQL migrations untuk membuat schema:
```bash
psql -U postgres -d pwa_db -f migrations/001_init_schema.sql
```

**Step 3: Run Application**

```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi
├── internal/
│   ├── config/                     # Konfigurasi aplikasi
│   ├── database/                   # Database connection & setup
│   ├── handlers/                   # HTTP handlers
│   ├── middleware/                 # Authentication middleware
│   ├── models/                     # Data models
│   └── repositories/               # Data access layer
├── migrations/                     # Database migrations
├── docs/                           # API documentation (Swagger)
├── Dockerfile                      # Docker configuration
├── docker-compose.yml              # Docker Compose configuration
├── go.mod & go.sum                 # Go dependencies
└── README.md                        # File ini
```

## API Documentation

API documentation tersedia dalam format Swagger/OpenAPI:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Swagger JSON**: `http://localhost:8080/swagger/swagger.json`
- **Swagger YAML**: `http://localhost:8080/swagger/swagger.yaml`

## Available Endpoints

Endpoint utama:
- **Auth**: `/api/auth/*` - Authentication & user management
- **Products**: `/api/products/*` - Product management
- **Transactions**: `/api/transactions/*` - Transaction management

Lihat Swagger documentation untuk detail endpoint lengkap.

## Development

### Run dengan Hot Reload (Recommended)

Install `air` untuk auto-reload:
```bash
go install github.com/cosmtrek/air@latest
```

Jalankan dengan:
```bash
air
```

### Build Binary

```bash
go build -o pwa-backend cmd/api/main.go
./pwa-backend
```

### Run Tests

```bash
go test ./...
```

## Troubleshooting

### Database Connection Error
- Pastikan PostgreSQL berjalan
- Cek DATABASE_URL di `.env` sudah benar
- Jika menggunakan Docker, pastikan network sudah terhubung dengan benar

### Port Already in Use
- Ubah PORT di `.env` ke port yang tersedia
- Atau kill process yang menggunakan port 8080:
  ```bash
  lsof -ti:8080 | xargs kill -9
  ```

### Migration Error
- Pastikan file migration ada di folder `migrations/`
- Cek permission untuk execute SQL files

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `JWT_SECRET` | Secret key untuk JWT token | - |
| `PORT` | Server port | 8080 |

## License

Project ini adalah bagian dari PWA Backend.
