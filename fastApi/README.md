# FastAPI RAG (Retrieval-Augmented Generation)

Proyek ini merupakan implementasi **Retrieval-Augmented Generation (RAG)** menggunakan **FastAPI** sebagai backend API dan **Large Language Model (LLM)** untuk menghasilkan jawaban berbasis dokumen yang relevan.

RAG mengombinasikan **pencarian vektor (retrieval)** dengan **model bahasa (generation)** sehingga jawaban yang dihasilkan lebih akurat, kontekstual, dan dapat dipertanggungjawabkan.

---

## 🚀 Fitur Utama
- FastAPI sebagai REST API backend
- Embedding teks untuk pencarian semantik
- Integrasi dengan Vector Database (misalnya Pinecone)
- Integrasi LLM (OpenAI)
- Arsitektur modular & mudah dikembangkan
- Dukungan environment variable (`.env`)

---

## 🏗️ Arsitektur Singkat
1. **User Query** dikirim ke API
2. Query diubah menjadi **embedding**
3. Embedding dicocokkan dengan **vector database**
4. Dokumen relevan dikirim ke **LLM**
5. LLM menghasilkan **jawaban kontekstual**

---

## 📂 Struktur Folder
FastAPI_RAG/
├── main.py # Entry point FastAPI
├── rag.py # Logika utama RAG
├── embedding.py # Proses embedding teks
├── llm.py # Konfigurasi dan pemanggilan LLM
├── pinecone_client.py # Koneksi ke Pinecone
├── requirements.py # Daftar dependency Python
├── README.md
├── .env.example # Contoh environment variable
└── .gitignore

yaml
Salin kode

---

## ⚙️ Instalasi & Setup

### 1️⃣ Clone Repository
```bash
git clone https://github.com/marioklau/FastAPI_RAG.git
cd FastAPI_RAG
2️⃣ Buat Virtual Environment
bash
Salin kode
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
.venv\Scripts\activate     # Windows
3️⃣ Install Dependency
bash
Salin kode
pip install -r requirements.txt
4️⃣ Konfigurasi Environment Variable
Salin file .env.example menjadi .env:

bash
Salin kode
cp .env.example .env
Isi .env:

env
Salin kode
OPENAI_API_KEY=your_openai_api_key
PINECONE_API_KEY=your_pinecone_api_key
PINECONE_ENV=your_pinecone_environment
⚠️ Jangan pernah mengunggah file .env ke repository publik

▶️ Menjalankan Aplikasi
bash
Salin kode
uvicorn main:app --reload
Akses API di:

cpp
Salin kode
http://127.0.0.1:8000
Dokumentasi otomatis:

Swagger UI: http://127.0.0.1:8000/docs

ReDoc: http://127.0.0.1:8000/redoc

🧪 Contoh Use Case
Chatbot berbasis dokumen internal

Sistem tanya jawab akademik

Knowledge base perusahaan

Asisten AI berbasis data lokal

🔐 Keamanan
API key disimpan menggunakan environment variable

.env diabaikan oleh Git (.gitignore)

Disarankan melakukan key rotation secara berkala

🛠️ Teknologi yang Digunakan
Python

FastAPI

OpenAI API

Pinecone Vector Database

Uvicorn

📌 Catatan Pengembangan
Proyek ini masih dapat dikembangkan dengan:

Autentikasi & rate limiting

Logging & monitoring

Dockerisasi

CI/CD (GitHub Actions)

Frontend (Next.js / React)

👤 Author
Mario Klau

📄 Lisensi
Proyek ini dibuat untuk keperluan pembelajaran dan pengembangan.
Silakan digunakan dan dimodifikasi sesuai kebutuhan.

yaml
Salin kode

---

## ✅ Langkah Setelah Ini
1. **Save `README.md`**
2. Jalankan:
```powershell
git add README.md
git commit -m "Add project README"
git push
