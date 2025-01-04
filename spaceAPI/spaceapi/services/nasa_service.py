import requests
from spaceapi.config.nasa_connection import Config
from spaceapi.utils.exceptions import APIError

class NasaService:
    @staticmethod
    def get_apod():
        try:
            response = requests.get(Config.NASA_BASE_URL, params={"api_key": Config.NASA_API_KEY})
            data = response.json()

            if response.status_code != 200:
                raise APIError(response.status_code, data.get("msg", "Error fetching data"))

            return data
        except requests.RequestException as e:
            raise APIError(500, f"Connection error: {str(e)}")