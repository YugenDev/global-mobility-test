from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from spaceapi.routes import apod

app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": " from space! 🚀"}

app.include_router(apod.router)
