from fastapi import APIRouter, Request, HTTPException
from fastapi.templating import Jinja2Templates
from spaceapi.services.nasa_service import NasaService
from spaceapi.utils.exceptions import APIError

router = APIRouter()
templates = Jinja2Templates(directory="spaceapi/templates")

@router.get("/apod")
def get_apod():
    try:
        data = NasaService.get_apod()
        return data
    except APIError as e:
        raise HTTPException(status_code=e.status_code, detail=e.message)


@router.get("/apod/html")
def get_apod_html(request: Request):
    try:
        data = NasaService.get_apod()
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
    except APIError as e:
        raise HTTPException(status_code=e.status_code, detail=e.message)