from typing import Any

from pydantic import BaseModel


class APIResponse(BaseModel):
    statusCode: int
    message: str
    data: Any = None