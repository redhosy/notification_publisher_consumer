# Notification Publisher Consumer System

Sistem notifikasi berbasis RabbitMQ yang diimplementasikan dengan Golang. Proyek ini terdiri dari dua komponen utama: publisher dan consumer, yang berkomunikasi melalui message broker RabbitMQ.

## Deskripsi Proyek

Proyek ini merupakan contoh implementasi pola publish-subscribe menggunakan RabbitMQ sebagai message broker. Dengan arsitektur ini, komponen publisher dapat mengirimkan pesan notifikasi ke dalam antrian (queue), dan komponen consumer dapat memproses pesan tersebut secara asinkron.

### Fitur Utama

- **Publisher**: Aplikasi yang dapat mengirim pesan notifikasi baik secara manual (melalui input pengguna) maupun otomatis (setiap 10 detik).
- **Consumer**: Aplikasi yang menerima dan memproses pesan dari antrian RabbitMQ.
- **Paket RabbitMQ**: Abstraksi untuk interaksi dengan RabbitMQ, termasuk penanganan koneksi, channel, exchange, dan queue.

## Persyaratan

- Go 1.22 atau lebih baru
- RabbitMQ Server (berjalan di localhost:5672 atau dapat dikonfigurasi)
- Make (untuk menjalankan perintah Makefile)

## Instalasi

1. Clone repositori ini:
   ```
   git clone https://github.com/redhosy/notification_publisher_consumer.git
   cd notification_publisher_consumer
   ```

2. Install dependensi:
   ```
   go mod tidy
   ```

3. Build aplikasi:
   ```
   make build
   ```

## Penggunaan

### Menjalankan RabbitMQ

Pastikan RabbitMQ server telah berjalan sebelum menjalankan publisher atau consumer. Jika Anda belum memiliki RabbitMQ, Anda dapat menjalankannya dengan Docker:

```
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

### Menjalankan Consumer

```
make run-consumer
```

atau

```
go run cmd/consumer/main.go
```

### Menjalankan Publisher

```
make run-publisher
```

atau

```
go run cmd/publisher/main.go
```

Setelah menjalankan publisher, Anda dapat:
1. Mengirim pesan manual dengan memasukkan teks dan menekan Enter
2. Ketik "exit" untuk keluar dari aplikasi
3. Pesan otomatis juga akan dikirim setiap 10 detik

## Struktur Proyek

```
.
├── cmd
│   ├── consumer         # Aplikasi consumer
│   │   └── main.go
│   └── publisher        # Aplikasi publisher
│       └── main.go
├── pkg
│   └── rabbitmq         # Paket untuk interaksi dengan RabbitMQ
│       └── rabbitmq.go
├── bin                  # Binary hasil kompilasi (dibuat setelah build)
├── go.mod               # Dependensi Go
├── go.sum
├── Makefile             # Perintah untuk build dan run
└── README.md
```

## Cara Kerja

1. **Publisher**:
   - Membuat koneksi ke RabbitMQ
   - Mendeklarasikan exchange dan queue
   - Menerima input dari pengguna atau membuat pesan otomatis
   - Mempublikasikan pesan ke exchange dengan routing key tertentu

2. **Consumer**:
   - Membuat koneksi ke RabbitMQ
   - Mendeklarasikan exchange dan queue yang sama dengan publisher
   - Mendaftarkan consumer untuk menerima pesan
   - Memproses pesan yang diterima dan mengirim acknowledgement

## Konfigurasi

Konfigurasi default:
- RabbitMQ host: localhost:5672
- Username/password: guest/guest
- Queue name: notification_queue
- Exchange: notification_exchange
- Routing key: notification_key

Untuk mengubah konfigurasi koneksi RabbitMQ, Anda dapat memodifikasi parameter URI pada saat membuat instance RabbitMQ di file main.go.

## Pengembangan Lebih Lanjut

Beberapa ide untuk pengembangan lebih lanjut:
- Menambahkan konfigurasi melalui file .env atau command line flags
- Implementasi retry mechanism untuk pesan yang gagal diproses
- Menambahkan logging dan monitoring
- Mengimplementasikan pola message routing yang lebih kompleks
- Menambahkan validasi dan skema pesan

## Lisensi

[MIT License](LICENSE)

## Kontributor

- [Redho Septa Yudien](https://github.com/redhosy)
