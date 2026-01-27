from dotenv import load_dotenv
load_dotenv()

from pinecone import Pinecone
from langchain_pinecone import PineconeVectorStore
from langchain_openai import OpenAIEmbeddings
import os

pc = Pinecone(api_key=os.environ["PINECONE_API_KEY"])
print("PINECONE:", os.getenv("PINECONE_API_KEY"))

# model embedding HARUS sama dengan saat indexing
embeddings = OpenAIEmbeddings(model="text-embedding-3-small")

# konek ke index 'exam' yang sudah kamu isi dari Colab
vectorstore = PineconeVectorStore.from_existing_index(
    index_name=os.environ["PINECONE_INDEX"],  # exam
    embedding=embeddings
)

# ini yang akan dipakai retriever
retriever = vectorstore.as_retriever(
    search_type="similarity_score_threshold",
    search_kwargs={
        "k": 6,
        "score_threshold": 0.75
    }
)

def invoke_by_material(material_id: int, query: str = ""):
    # Menggunakan similarity_search langsung agar filter berfungsi 100%
    # Filter ini akan memastikan HANYA vector dengan metadata `material_id` yang sama yang diambil.
    docs = vectorstore.similarity_search(
        query=query,
        k=6,
        filter={"material_id": material_id}
    )
    return docs
