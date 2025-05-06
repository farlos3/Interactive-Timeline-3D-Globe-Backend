from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware

import json

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "Python Service is running"}

@app.post("/process")
async def process_data(request: Request):
    body = await request.body()
    if not body:
        print("=== Body is empty! ===")
        return JSONResponse(content={"status": "error", "message": "Body is empty"}, status_code=400)
    data = await request.json()
    print("=== Data received from GO ===")
    print(json.dumps(data, indent=4, ensure_ascii=False))
    return JSONResponse(content={"status": "success", "received": data})

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)