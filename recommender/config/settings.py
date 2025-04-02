import os
from dotenv import load_dotenv

load_dotenv()

PRODUCT_API = os.getenv("PRODUCT_API")
KAFKA_SERVER = os.getenv("KAFKA_SERVER")
KAFKA_PORT = os.getenv("KAFKA_PORT")
