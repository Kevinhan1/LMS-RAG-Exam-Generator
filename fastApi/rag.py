from pinecone_client import retriever
from langchain_openai import ChatOpenAI
from langchain_core.prompts import PromptTemplate

EXAM_PROMPT_STRICT = PromptTemplate(
    input_variables=["context", "instruction"],
    template="""
Anda adalah AI pembuat soal ujian PILIHAN GANDA.

⚠️ ATURAN KERAS:
- GUNAKAN HANYA informasi dari KONTEKS
- DILARANG menggunakan pengetahuan umum
- Jika konteks TIDAK MEMUAT materi yang diminta → KEMBALIKAN ERROR (BUKAN SOAL)

### KONTEKS
{context}
### AKHIR KONTEKS

### PERMINTAAN atau INSTRUKSI TEACHER
{instruction}

### FORMAT OUTPUT (WAJIB JSON VALID)
- HANYA kembalikan JSON
- Output WAJIB berupa ARRAY/LIST dari object soal, contoh: `[ {{...}}, {{...}} ]`
- Jika diminta 1 soal, tetap kembalikan dalam array `[ {{...}} ]`
- Gunakan format PERSIS seperti berikut:

[
  {
    "material_id": {material_id},
    "content": "Teks soal murni tanpa nomor...",
    "difficulty": "easy | medium | hard",
    "taxonomy_level": "C5 | C6",
    "answers": [
      { "label": "A", "text": "...", "is_correct": false },
      { "label": "B", "text": "...", "is_correct": true }
    ]
  },
  ...
]

⚠️ ATURAN TAMBAHAN:
- content WAJIB MURNI TEKS SOAL. JANGAN pakai "Pertanyaan 1", "No. 1", dll.
- Jumlah soal WAJIB SESUAI instruksi teacher (jika tidak disebut, buat 1).
- Taxonomy Level & Difficulty:
  1. IKUTI instruksi teacher jika ada.
  2. JIKA TIDAK ADA instruksi, TENTUKAN SENDIRI berdasarkan analisis soal.
  3. Taxonomy Level WAJIB HANYA "C5" atau "C6". JANGAN gunakan C1-C4.
  4. Difficulty soal bebas, jika tidak ada instruksi, tentukan sendiri berdasarkan analisis soal: "easy", "medium", atau "hard".

- TIDAK BOLEH ada penjelasan / pembahasan
- TIDAK BOLEH mengulang konteks
- TEPAT SATU jawaban is_correct = true
- Jumlah pilihan jawaban MENYESUAIKAN INSTRUKSI (default 4 jika tidak disebut)
- Label berurutan A, B, C, D...
"""
)

llm = ChatOpenAI(
    model="gpt-4o-mini-2024-07-18",
    temperature=0.3
)

def validate_context(docs):
    if not docs:
        raise ValueError("Materi tidak ditemukan di Pinecone.")

    text = " ".join(doc.page_content for doc in docs)
    if len(text) < 300:
        raise ValueError("Konteks terlalu sedikit untuk generate soal.")

    return text


def generate_exam(material_id: int, instruction: str) -> dict:
    """
    FULL RAG PIPELINE
    """

    # 1️⃣ Ambil context dari Pinecone BERDASARKAN material_id
    docs = retriever.invoke_by_material(
        material_id=material_id,
        query=instruction
    )

    # 2️⃣ Validasi context
    context = validate_context(docs)

    # 3️⃣ Build prompt
    prompt = EXAM_PROMPT_STRICT.format(
        context=context,
        instruction=instruction,
        material_id=material_id
    )

    # 4️⃣ Call LLM
    response = llm.invoke(prompt)

    return response.content
