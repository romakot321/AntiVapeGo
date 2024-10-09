import aiohttp
import random
import asyncio
import time

baseurl = "http://localhost:8000"
guids = ["string", "killer", "a", "b"]


async def create_sensor_data(session):
    data = {
        "guid": guids[int(time.time() * 100 % 4)],
        "co2": random.randint(1, 2000),
        "tvoc": random.randint(1, 2000),
        "batteryCharge": random.randint(1, 120)
    }
    async with session.post("/external/sensors_data", json=data) as resp:
        assert resp.status == 200


async def main(total: int):
    tasks = []
    async with aiohttp.ClientSession(baseurl) as session:
        for _ in range(total):
            if total % (total / 100) == 0:
                print(total / 100)
            tasks.append(create_sensor_data(session))
        await asyncio.gather(*tasks)
    print()


asyncio.run(main(100000))
