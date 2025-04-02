from sqlalchemy import Column, String, Float, Integer, DateTime, func
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class Product(Base):
    __tablename__ = "products"
    id = Column(String, primary_key=True)
    name = Column(String)
    description = Column(String)
    price = Column(Float)
    account_id = Column(Integer)

class Interaction(Base):
    __tablename__ = "interactions"
    id = Column(Integer, primary_key=True, autoincrement=True)
    user_id = Column(String)
    product_id = Column(String)
    interaction_type = Column(String)
    timestamp = Column(DateTime, default=func.now())
