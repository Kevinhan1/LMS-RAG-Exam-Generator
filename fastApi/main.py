import requests
import json
import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from fastapi import UploadFile, File, Form
from ingest import ingest_pdf
from rag import generate_exam

app = FastAPI()

BACKEND_URL = os.getenv("BACKEND_URL")

class ExamRequest(BaseModel):
    material_id: int
    instruction: str

@app.post("/generate_exam")
def generate(data: ExamRequest):
    try:
        # 1️⃣ Jalankan RAG (SEMUA logic di rag.py)
        result = generate_exam(
            material_id=data.material_id,
            instruction=data.instruction
        )

        # 2️⃣ Parse JSON dari LLM
        exam_json = json.loads(result)

        return {
            "rag_result": exam_json
        }

    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

    except json.JSONDecodeError:
        raise HTTPException(status_code=500, detail="LLM output bukan JSON valid")

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/ingest_material")
async def ingest_material(
    file: UploadFile = File(...),
    material_id: int = Form(...),
    course_id: int = Form(...),
    chapter_id: int = Form(...)
):
    result = await ingest_pdf(
        file,
        material_id,
        course_id,
        chapter_id
    )
    return result