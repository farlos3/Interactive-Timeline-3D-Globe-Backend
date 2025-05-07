from fastapi import FastAPI
from package.router import router
from package.middleware import middlewares

app = FastAPI(
    title="War Events Clustering API",
    description="API for clustering war events using spatial and temporal data",
    version="1.0.0"
)

middlewares(app)

app.include_router(router)