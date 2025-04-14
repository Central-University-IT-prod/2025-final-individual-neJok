import aiohttp
from config import API_URL

async def get_client_data(client_id: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/clients/{client_id}") as response:
            if response.status == 200:
                return await response.json()
            else:
                return None
            
async def get_client_ads(client_id: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/ads?client_id={client_id}") as response:
            if response.status == 200:
                return await response.json()
            else:
                return None

async def get_current_date() -> int | None:
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/time/advance") as response:
            if response.status == 200:
                return (await response.json())['current_date']
            else:
                return None

async def ad_click(client_id: str, ad_id: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        data = {
            "client_id": client_id,
        }
        async with session.post(f"{API_URL}/ads/{ad_id}/click", json=data) as response:
            if response.status == 204:
                return True
            else:
                return None

async def get_advertiser_data(advertiser_id: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/advertisers/{advertiser_id}") as response:
            if response.status == 200:
                return await response.json()
            else:
                return None
            

async def get_advertiser_campaigns(advertiser_id: str, page: int, size: int = 2):
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/advertisers/{advertiser_id}/campaigns?page={page}&size={size}") as response:
            if response.status == 200:
                count_campaigns = response.headers.getone("X-Total-Count")
                return int(count_campaigns), await response.json()
            else:
                return 0, None

async def get_campaign(advertiser_id: str, campaign_id: str):
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/advertisers/{advertiser_id}/campaigns/{campaign_id}") as response:
            if response.status == 200:
                return await response.json()
            else:
                return None
            
async def create_client(client_id: str, login: str, age: str, location: str, gender: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        url = f"{API_URL}/clients/bulk"
        data = {
            "client_id": client_id,
            "login": login,
            "age": age,
            "location": location,
            "gender": gender
        }
        async with session.post(url, json=[data]) as response:
            if response.status == 201:
                return await response.json()
            else:
                return None

            
async def create_advertiser(advertiser_id: str, name: str) -> dict | None:
    async with aiohttp.ClientSession() as session:
        url = f"{API_URL}/advertisers/bulk"
        data = {
            "advertiser_id": advertiser_id,
            "name": name
        }
        async with session.post(url, json=[data]) as response:
            if response.status == 201:
                return await response.json()
            else:
                return None
            
async def create_campaign(
    advertiser_id: str,
    impressions_limit: int,
    clicks_limit: int,
    cost_per_impression: float,
    cost_per_click: float,
    ad_title: str,
    ad_text: str,
    start_date: int,
    end_date: int,
    targeting: dict
) -> dict | None:
    async with aiohttp.ClientSession() as session:
        url = f"{API_URL}/advertisers/{advertiser_id}/campaigns"
        data = {
            "impressions_limit": impressions_limit,
            "clicks_limit": clicks_limit,
            "cost_per_impression": cost_per_impression,
            "cost_per_click": cost_per_click,
            "ad_title": ad_title,
            "ad_text": ad_text,
            "start_date": start_date,
            "end_date": end_date,
            "targeting": targeting
        }
        async with session.post(url, json=data) as response:
            if response.status == 201:
                return await response.json()
            else:
                return None
