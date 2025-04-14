import hashlib
import uuid

def generate_uuid_from_id(telegram_id: int) -> str:
    telegram_id_str = str(telegram_id)
    hash_object = hashlib.sha256(telegram_id_str.encode())
    hash_bytes = hash_object.digest()[:16]
    generated_uuid = uuid.UUID(bytes=hash_bytes)
    
    return str(generated_uuid)