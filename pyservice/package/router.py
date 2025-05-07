# router.py
from fastapi import APIRouter, Request, HTTPException
from fastapi.responses import JSONResponse
from .calculation_service import CalculationService
import json
from typing import List, Dict
from datetime import datetime

router = APIRouter()
calculation_service = CalculationService()

@router.get("/")
async def root():
    return {"message": "Python Service is running"}

@router.post("/process")
async def process_data(request: Request):
    try:
        body = await request.body()
        if not body:
            print("=== Body is empty! ===")
            return JSONResponse(
                content={"status": "error", "message": "Body is empty"}, 
                status_code=400
            )

        data = await request.json()
        # print("=== Data received from GO ===")
        # print(json.dumps(data, indent=4, ensure_ascii=False))

        # ตรวจสอบรูปแบบข้อมูล
        if not isinstance(data, dict) or 'events' not in data:
            return JSONResponse(
                content={"status": "error", "message": "Invalid data format. Expected {'events': [...]}"}, 
                status_code=400
            )

        # ตรวจสอบข้อมูลแต่ละ event
        events = data['events']
        for idx, event in enumerate(events):
            # ตรวจสอบฟิลด์ที่จำเป็น
            required_fields = ['EventID', 'Lat', 'Lon', 'Date']
            missing_fields = [field for field in required_fields if field not in event]
            if missing_fields:
                return JSONResponse(
                    content={
                        "status": "error", 
                        "message": f"Missing required fields in event {idx}: {', '.join(missing_fields)}"
                    }, 
                    status_code=400
                )
            
            # ตรวจสอบประเภทข้อมูล
            try:
                float(event['Lat']) 
                float(event['Lon']) 
                
                # Preprocess วันที่
                date_str = event['Date']
                if 'T' in date_str:
                    date_str = date_str.split('T')[0] 
                event['Date'] = date_str
                
                datetime.strptime(date_str, "%Y-%m-%d")  # ตรวจสอบรูปแบบวันที่
            except ValueError as e:
                return JSONResponse(
                    content={
                        "status": "error",
                        "message": f"Invalid data type in event {idx}: {str(e)}"
                    },
                    status_code=400
                )

        # ประมวลผลข้อมูล
        result = calculation_service.process_events(events)
        
        return JSONResponse(
            content={
                "status": "success",
                "data": result
            }
        )

    except Exception as e:
        print(f"=== Error processing data: {str(e)} ===")
        return JSONResponse(
            content={"status": "error", "message": str(e)}, 
            status_code=500
        )