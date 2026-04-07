"""Example tool — delete or replace with real business tools."""

from pydantic import BaseModel

from .registry import ToolContext, tool


class EchoInput(BaseModel):
    message: str


class EchoOutput(BaseModel):
    echo: str


@tool(description="Echo the input message back")
async def echo(params: EchoInput, ctx: ToolContext) -> EchoOutput:
    return EchoOutput(echo=params.message)
