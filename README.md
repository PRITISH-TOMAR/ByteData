# ByteData

A high-performance data management application built with Go.

## **Quick Start**

### Option 1: Run with Docker Compose (Recommended)

#### 1. **Clone the repository**:
```bash
git clone https://github.com/PRITISH-TOMAR/bytedata.git
cd bytedata
```

#### 2. **Start the services in background**:
```bash
docker-compose up -d --build
```

#### 4. **Run the client interactively**:
```bash
# With specific username:
docker-compose run --rm client ./client -u <username>
```

#### 5. **Stop everything**:
```bash
docker-compose down
```
--- 
---

### Option 2: Run Directly (Manual Build)

#### 1. **Clone the repository**:
```bash
git clone https://github.com/PRITISH-TOMAR/bytedata.git
cd bytedata
```

#### 2. **Build and run the Server**:
```bash
cd Server
go mod download
go build -o server .
./server -port 9090
```

Keep this terminal running.

#### 3. **Build and run the Client** (in a new terminal):
```bash
cd Client
go mod download
go build -o client .

./client -u <username> -addr localhost:9090
```

---

## **Features**

- Database Server with WAL (Write-Ahead Logging) Persistence
- Multiple bucket support
- Range queries support over the bucket
- Multiple client connections over TCP


## **Requirements**

- **Option 1**: Docker and Docker Compose
- **Option 2**: Go 1.24 or higher



