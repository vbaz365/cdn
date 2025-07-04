# 🌐 CDN DNS Resolver

This project is a Go implementation of an **Authoritative DNS Resolver** that uses **EDNS Client Subnet (ECS)** information to return geographically optimized responses. It determines the closest CDN **Point of Presence (PoP)** for a given IPv6 prefix using a **compressed binary radix tree**.

---

## 📌 Problem Description

When a DNS query includes an ECS field, the server must:

1. Match the ECS IPv6 prefix to the most specific subnet in its internal **routing table**.
2. Return:
   - The **PoP ID** (associated with the matched subnet).
   - The **scope prefix-length** (length of the matched prefix).

For example:
Routing Entry: 2001:49f0:d0b8::/48 => PoP 174
ECS Query: 2001:49f0:d0b8:8a00::/56
Response: PoP 174, Scope 48

## 📦 Project Structure
```text
.
├── main.go                   # Loads data and performs ECS lookup from CLI
├── data/
│   └── routing-data.txt      # Input file with IPv6 prefixes and PoP IDs
├── datastructure/
│   └── radix_tree.go         # Compressed radix tree implementation
├── utils/
│   ├── radix_tree_utils.go   # Bit operations & prefix manipulation
│   └── radix_tree_utils_test.go # Unit tests for utils
```

---

## ⚙️ How It Works

The core of the project is a **compressed binary radix tree** that:

- Stores IPv6 subnet prefixes efficiently.
- Compresses paths by grouping consecutive matching bits.
- Traverses the tree to perform a **longest prefix match** on incoming ECS queries.

### 🔧 Routing Entry Format

Each line of `routing-data.txt` contains:
<IPv6 Prefix> <PoP ID>
Example:
2001:49f0:d0b8::/48 174

---

## 🚀 Getting Started

### 📥 Installation

```bash
git clone https://github.com/yourusername/cdn-dns-resolver
cd cdn-dns-resolver
go build -o resolver main.go
```

### 🧪 Running the Resolver
- ./resolver 2001:db8:abcd:1::/64

### 📤 Example Output
- Pop id: 200, Scope: 48

### 🧪 Running Unit Tests
To run the unit tests for utility functions:
- go test ./utils



