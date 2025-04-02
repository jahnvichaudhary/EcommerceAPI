from fastapi import APIRouter
from kafka import KafkaProducer
router = APIRouter(
    prefix="/v1",
    tags=["recommendation"],
)

@router.get("/recommend")
async def recommend(
        user_id: int,
):
    pass

