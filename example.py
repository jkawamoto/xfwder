# /// script
# dependencies = [
#   "fastapi>=0.115",
#   "uvicorn>=0.34",
# ]
# ///
import asyncio
import subprocess
import sys
import webbrowser
from asyncio import sleep
from http import HTTPStatus

from fastapi import FastAPI, Request
from starlette.datastructures import QueryParams
from uvicorn import Server, Config

UDS_FILENAME = "e2e"


async def main() -> None:
    res = subprocess.call(["open", "build/XFwder.app", "--args", "--init"])
    if res != 0:
        sys.exit(res)
    await sleep(5)

    api = FastAPI()
    server = Server(Config(api, uds=f"/tmp/{UDS_FILENAME}.sock"))

    sample_queries = QueryParams({"arg1": "value1", "arg2": "value2"})

    @api.post("/e1", status_code=HTTPStatus.NO_CONTENT)
    def e1(req: Request) -> None:
        assert req.query_params == sample_queries
        server.should_exit = True

    server_task = asyncio.create_task(server.serve())

    webbrowser.open(f"xfwder://{UDS_FILENAME}/e1?{sample_queries}")

    await server_task


if __name__ == "__main__":
    asyncio.run(main())
