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
    # Register the custom URL scheme xfwder:// by launching XFwder.app with --init flag
    # This only needs to be done once during initial setup
    res = subprocess.call(["open", "build/XFwder.app", "--args", "--init"])
    if res != 0:
        sys.exit(res)
    await sleep(5)  # Wait for app initialization to complete

    # Define query parameters that will be passed in the URL
    sample_queries = QueryParams({"arg1": "value1", "arg2": "value2"})

    # Prepare the receiver server that listens to on a UNIX domain socket.
    api = FastAPI()
    server = Server(Config(api, uds=f"/tmp/{UDS_FILENAME}.sock"))

    @api.post("/e1", status_code=HTTPStatus.NO_CONTENT)
    def e1(req: Request) -> None:
        # Verify that the received request contains expected query parameters
        assert req.query_params == sample_queries
        # Close the server to exit this example script.
        server.should_exit = True

    server_task = asyncio.create_task(server.serve())

    # Trigger the request by opening URL with the custom scheme: xfwder://{socket_name}/endpoint?params
    webbrowser.open(f"xfwder://{UDS_FILENAME}/e1?{sample_queries}")

    # Wait for the server to process the request and shut down
    await server_task


if __name__ == "__main__":
    asyncio.run(main())
