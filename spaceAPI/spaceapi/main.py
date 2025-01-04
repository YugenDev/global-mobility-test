from fastapi import FastAPI, HTTPException, Request
from fastapi.templating import Jinja2Templates
import requests
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()
templates = Jinja2Templates(directory="spaceapi/templates")

NASA_API_KEY = os.getenv("NASA_API_KEY")

@app.get("/")
def read_root():
    return {"Hello": " from space! 🚀"}

@app.get("/apod")
def get_apod():
    try:
        url = "https://api.nasa.gov/planetary/apod"
        params = {"api_key": NASA_API_KEY}

        response = requests.get(url, params=params)
        data = response.json()

        if response.status_code != 200:
            raise HTTPException(status_code=response.status_code, detail=data.get("msg", "error fetching data"))

        return {
            "title": data["title"],
            "date": data["date"],
            "explanation": data["explanation"],
            "url": data["url"],
            "media_type": data["media_type"]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    
@app.get("/apod/html")
def get_apod_html(request: Request):
    try:
        url = "https://api.nasa.gov/planetary/apod"
        params = {"api_key": NASA_API_KEY}

        response = requests.get(url, params=params)
        data = response.json()

        if response.status_code != 200:
            raise HTTPException(status_code=response.status_code, detail=data.get("msg", "error fetching data"))
        return templates.TemplateResponse(
            "apod.html",
            {
                "request": request,
                "title": data["title"],
                "date": data["date"],
                "explanation": data["explanation"],
                "url": data["url"],
                "media_type": data["media_type"]
            }
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
