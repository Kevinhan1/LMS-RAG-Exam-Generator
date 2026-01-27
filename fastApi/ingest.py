from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.documents import Document
from pinecone_client import vectorstore
import tempfile
import os
from pypdf import PdfReader

async def ingest_pdf(file, material_id, course_id, chapter_id):
    # simpan sementara
    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as tmp:
        tmp.write(await file.read())
        tmp_path = tmp.name

    reader = PdfReader(tmp_path)
    text = ""
    for page in reader.pages:
        text += page.extract_text() + "\n"

    os.unlink(tmp_path)

    splitter = RecursiveCharacterTextSplitter(
        chunk_size=800,
        chunk_overlap=150
    )

    chunks = splitter.split_text(text)

    docs = [
        Document(
            page_content=chunk,
            metadata={
                "material_id": material_id,
                "course_id": course_id,
                "chapter_id": chapter_id,
                "source": file.filename
            }
        )
        for chunk in chunks
    ]

    vectorstore.add_documents(docs)

    return {
        "status": "ok",
        "chunks": len(docs)
    }
