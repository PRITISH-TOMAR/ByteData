# ByteData

A high-performance data management application built with Go.

## Quick Start

### Option A: Run Directly (from Git)

1. **Clone and run**:
```bash
git clone https://github.com/PRITISH-TOMAR/bytedata.git
cd bytedata
go run cmd/main.go -u root -p root
```

2. **Or build and run**:
```bash
git clone https://github.com/PRITISH-TOMAR/bytedata.git
cd bytedata
go build -o bytedata cmd/main.go
./bytedata -u root -p root
```

The application will run with default configuration.

### Option B: Run with Docker (with volume mount)

1. **Clone the repository**:
```bash
git clone https://github.com/PRITISH-TOMAR/bytedata.git
cd bytedata
```

2. **Build and run with Docker**:
```bash
# Build the image
docker build -t bytedata .

# For the very first time, Run with volume mount for data persistence 
docker run -it  -p 4040:4040 -v ./.bytedata:/root/.bytedata/ bytedata

# After setting up the new username and password
docker run -it  -p 4040:4040 -v $(pwd)/.bytedata:/root/.bytedata/ bytedata -u \<username\>
```
