from surprise import SVD, Dataset, Reader
import pandas as pd
from app.db.session import ReplicaSession
from app.db.models import Interaction, Product

def fetch_interactions() -> pd.DataFrame:
    with ReplicaSession() as session:
        interactions = session.query(Interaction).all()
        data = [
            {
                "user_id": i.user_id,
                "product_id": i.product_id,
                "rating": 3.0 if i.interaction_type == "purchase" else 1.0
            }
            for i in interactions
        ]
        return pd.DataFrame(data)

class Recommender:
    def __init__(self):
        self.model = SVD(n_factors=50, random_state=42)
        self.trainset = None
        self.product_ids = set()

    def train(self):
        df = fetch_interactions()
        self.product_ids = set(df["product_id"].unique())
        reader = Reader(rating_scale=(1, 3))
        data = Dataset.load_from_df(df[["user_id", "product_id", "rating"]], reader)
        self.trainset = data.build_full_trainset()
        self.model.fit(self.trainset)

    def recommend(self, user_id: str, top_n: int = 5) -> list[str]:
        with ReplicaSession() as session:
            self.product_ids = {p.id for p in session.query(Product.id).all()}
            interacted = {
                i.product_id for i in session.query(Interaction.product_id)
                .filter(Interaction.user_id == user_id).all()
            }
        candidates = [pid for pid in self.product_ids if pid not in interacted]
        predictions = [self.model.predict(user_id, pid) for pid in candidates]
        top_predictions = sorted(predictions, key=lambda x: x.est, reverse=True)[:top_n]
        return [pred.iid for pred in top_predictions]

recommender = Recommender()  # Singleton for simplicity
