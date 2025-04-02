from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from app.db.models import Base

replica_engine = create_engine("")
ReplicaSession = sessionmaker(bind=replica_engine)
Base.metadata.create_all(replica_engine)

def get_db():
    db = ReplicaSession()
    try:
        yield db
    finally:
        db.close()