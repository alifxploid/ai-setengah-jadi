# 🔐 Key Generator Tool

Tool untuk generate JWT secret dan access key secara otomatis dengan keamanan tinggi.

## 🚀 Cara Penggunaan

### Menggunakan Makefile (Recommended)
```bash
make generate-keys
```

### Menjalankan Langsung
```bash
cd cmd/generate-keys
go run main.go
```

## ✨ Fitur

- **Cryptographically Secure**: Menggunakan `crypto/rand` untuk generate random bytes
- **Base64 URL Encoding**: Aman untuk digunakan di environment variables
- **Auto Update .env**: Otomatis update file `.env` dengan keys baru
- **Timestamp Tracking**: Mencatat kapan keys di-generate
- **Safe Overwrite**: Mengganti keys lama tanpa merusak konfigurasi lain

## 🔑 Keys yang Di-generate

### JWT_SECRET
- **Length**: 64 bytes (512 bits)
- **Encoding**: Base64 URL
- **Usage**: Untuk signing dan verifying JWT tokens
- **Security**: Minimum 32 karakter untuk keamanan optimal

### ACCESS_KEY
- **Length**: 32 bytes (256 bits) 
- **Encoding**: Base64 URL
- **Usage**: Untuk admin access dan API authentication
- **Security**: Cryptographically secure random generation

## 📁 File yang Dimodifikasi

- `../../.env` - File environment variables utama
- Menambahkan timestamp comment untuk tracking
- Preserve semua konfigurasi existing

## ⚠️ Security Notes

1. **Never commit keys to version control**
2. **Regenerate keys untuk production environment**
3. **Store keys securely di production**
4. **Rotate keys secara berkala**
5. **Backup keys sebelum regenerate**

## 🔄 Regenerating Keys

Untuk generate ulang keys:

```bash
# Backup current .env (optional)
cp .env .env.backup

# Generate new keys
make generate-keys
```

## 🛠️ Development

Tool ini menggunakan:
- `crypto/rand` untuk secure random generation
- `regexp` untuk pattern matching di .env file
- `base64` untuk safe encoding
- `time` untuk timestamp tracking

## 📝 Example Output

```
🔐 Generating secure keys...
✅ Keys generated and .env updated successfully!
🔑 JWT_SECRET: sEJmU2nN5D1oQrEm07i13UURnSXBIVIz42bHFz7pMAP_pxoSPEiJgXJgyAM0xXCcxNJx3xSSFeXyeGXAGUACLQ==
🗝️  ACCESS_KEY: NBzbgYqdG6oErJPOKzi4JkaFK3eka8C5TPcz4uLikuY=

⚠️  Keep these keys secure and never commit them to version control!
```