# 🐼 PandaFS  
<div align="center">
    <img src="https://github.com/user-attachments/assets/6f8dbe11-40c1-447e-8daf-0bf996df75dc" alt="Banner Image" width="100%">
</div>

*A GFS-like distributed file system for scalable and reliable storage.*  

PandaFS is a simple **distributed file system** based on the **Master-Slave architecture**, designed for **fault tolerance, replication, and security**. It ensures reliable data storage and retrieval across multiple nodes.  

---

##  Features  

- **Fault Tolerant** – Ensures data availability even if nodes fail.  
- **Replication** – Data is duplicated across nodes to guarantee eventual consistency.  
- **Secure** – 
    - Supports **client-side encryption** for secure storage 
    - **ACL** and  **token-based authentication** mechanism for unauthorized access.  
- **High Performance** – Supports parallel reads and writes for efficiency.  

---

##  Architecture  

PandaFS follows a **Master-Slave architecture**:  

- **Master Node**: Manages metadata, replication, and client requests.  
- **Slave Nodes**: Store actual data and serve client read/write operations.  
- **Clients**: Interact with the system using a simple gRPC API.  

---


## Using client-sdk

```
go get github.com/nanoDFS/client-sdk
```

```go
// upload

key := crypto.DefaultCryptoKey()

fileId, userId, err := fs.NewFileSystem().Upload(key, "./test.mp4")

// download

err = fs.NewFileSystem().Download(fileId, userId, key, "./temp")
```
