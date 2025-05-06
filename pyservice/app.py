from fastapi import FastAPI
from package.router import router
from package.middleware import middlewares

app = FastAPI()
middlewares(app)
app.include_router(router)