# router.py
from fastapi import APIRouter, Request
from fastapi.responses import JSONResponse
import json

router = APIRouter()

@router.get("/")
async def root():
    return {"message": "Python Service is running"}

@router.post("/process")
async def process_data(request: Request):
    body = await request.body()
    if not body:
        print("=== Body is empty! ===")
        return JSONResponse(content={"status": "error", "message": "Body is empty"}, status_code=400)
    data = await request.json()
    print("=== Data received from GO ===")
    print(json.dumps(data, indent=4, ensure_ascii=False))
    return JSONResponse(content={"status": "success", "received": data})